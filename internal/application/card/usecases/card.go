package usecases

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters"
	cardRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/notification/constants"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/notification/services"
	log "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/observability/adapters"
	userRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/repository"
	errorsBusiness "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/card/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/card/projections"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/card/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/card/dto/response"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/card/entities"
)

const (
	cardsPerBatch = 3 // Number of cards per batch
	numWorkers    = 5 // Number of concurrent workers
)

type CardUsecaseIF interface {
	CreateCard(
		ctx context.Context,
		request request.CardRequest,
		enterpriseId string,
		preferredLanguage string,
	) (response.CardResponse, error)

	TriggerCardExpiringSoonNotifications(
		ctx context.Context,
		request request.NotificationCardExpiringSoonRequestDTO,
	) (response.NotificationCardExpiringSoonResponseDTO, error)
}

type CardUseCase struct {
	cardWriteRepository cardRepository.CardWriteRepositoryIF
	cardReadRepository  cardRepository.CardReadRepositoryIF
	userReadRepository  userRepository.UserReadRepositoryIF
	cardAdapter         adapters.CardAdapter
	notificationService services.NotificationService
	log                 log.Logger
}

func NewCardUseCase(
	cardAdapter adapters.CardAdapter,
	cardWriteRepository cardRepository.CardWriteRepositoryIF,
	cardReadRepository cardRepository.CardReadRepositoryIF,
	userReadRepository userRepository.UserReadRepositoryIF,
	notificationService services.NotificationService,
	log log.Logger,
) CardUsecaseIF {
	return &CardUseCase{
		cardWriteRepository: cardWriteRepository,
		cardReadRepository:  cardReadRepository,
		userReadRepository:  userReadRepository,
		cardAdapter:         cardAdapter,
		notificationService: notificationService,
		log:                 log,
	}
}

func (p *CardUseCase) CreateCard(
	ctx context.Context,
	request request.CardRequest,
	enterpriseID string,
	preferredLanguage string,
) (response.CardResponse, error) {
	newCard := entities.NewCard(
		request,
		enterpriseID,
	)

	createCardErr := p.cardWriteRepository.CreateCard(ctx, &newCard)

	if createCardErr != nil {
		return response.CardResponse{}, errorsBusiness.NewCreateCardDBError(createCardErr)
	}

	email, getEmailByUserIDErr := p.userReadRepository.GetEmailByUserID(ctx, request.UserID, enterpriseID)

	if getEmailByUserIDErr != nil {
		p.log.Warn(ctx, fmt.Sprintf("%v", getEmailByUserIDErr))
	}

	formattedActionDate := newCard.CreatedAt.UTC().
		In(time.FixedZone("UTC-6", -6*60*60)).Format("2006-01-02 15:04")

	notifyCardAdditionErr := p.notificationService.NotifyCardAddition(
		ctx,
		[]constants.NotificationChannel{constants.EmailChannel},
		&request.UserID,
		&email,
		&formattedActionDate,
		&newCard.LastFour,
		preferredLanguage,
	)

	if notifyCardAdditionErr != nil {
		p.log.Warn(ctx, fmt.Sprintf("%v", notifyCardAdditionErr))
	}

	return response.CardResponse{
		ID:             newCard.ID.String(),
		CardTokenID:    newCard.ExternalCardID,
		Alias:          newCard.Alias,
		LastFour:       newCard.LastFour,
		Brand:          newCard.Brand,
		IsDefault:      newCard.IsDefault,
		IsRecurrent:    newCard.IsRecurrent,
		ExpirationDate: newCard.ExpirationDate.Format("01/06"),
	}, nil
}

func (p *CardUseCase) TriggerCardExpiringSoonNotifications(
	ctx context.Context,
	request request.NotificationCardExpiringSoonRequestDTO,
) (response.NotificationCardExpiringSoonResponseDTO, error) {
	return decorators.TraceDecorator(
		ctx,
		"CardUseCase.TriggerCardExpiringSoonNotifications",
		func(ctx context.Context, span decorators.Span) (response.NotificationCardExpiringSoonResponseDTO, error) {
			nextMonthCurrentDate := time.Now().AddDate(0, 1, 0)

			notificationCardExpiringSoonProjections, err := p.cardReadRepository.GetCardsExpiringSoon(
				ctx,
				nextMonthCurrentDate.Month(),
				nextMonthCurrentDate.Year(),
			)
			if err != nil {
				return response.NotificationCardExpiringSoonResponseDTO{}, err
			}

			if len(notificationCardExpiringSoonProjections) == 0 {
				return response.NotificationCardExpiringSoonResponseDTO{
					Message: "no cards found to trigger expiring soon notifications",
				}, nil
			}

			expiringSoonCardJobs := make(
				chan []projections.NotificationCardExpiringSoonProjection,
				(len(notificationCardExpiringSoonProjections)+cardsPerBatch-1)/cardsPerBatch,
			)
			errChan := make(chan error, len(notificationCardExpiringSoonProjections))
			var wg sync.WaitGroup

			for range make([]struct{}, numWorkers) {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for expiringSoonCardBatch := range expiringSoonCardJobs {
						p.processRetrievedCards(
							ctx,
							expiringSoonCardBatch,
							request.NotificationChannels,
							errChan,
						)
					}
				}()
			}

			for i := 0; i < len(notificationCardExpiringSoonProjections); i += cardsPerBatch {
				end := min(i+cardsPerBatch, len(notificationCardExpiringSoonProjections))
				expiringSoonCardJobs <- notificationCardExpiringSoonProjections[i:end]
			}
			close(expiringSoonCardJobs)

			wg.Wait()
			close(errChan)

			for err := range errChan {
				p.log.Warn(ctx, fmt.Sprintf("%v", err))
			}

			return response.NotificationCardExpiringSoonResponseDTO{
				Message: "card expiring soon notification has been triggered",
			}, nil
		})
}

func (p *CardUseCase) processRetrievedCards(
	ctx context.Context,
	notificationCardExpiringSoonProjections []projections.NotificationCardExpiringSoonProjection,
	notificationChannels []constants.NotificationChannel,
	errChan chan error,
) {
	for _, projection := range notificationCardExpiringSoonProjections {
		err := p.notificationService.NotifyCardExpiringSoon(
			ctx,
			notificationChannels,
			projection,
		)

		if err != nil {
			errChan <- err
			continue
		}

		p.log.Info(ctx, "the notification has been triggered successfully")
	}
}
