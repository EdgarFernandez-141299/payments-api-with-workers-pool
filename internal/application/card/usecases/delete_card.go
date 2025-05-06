package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters"
	cardRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/notification/constants"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/notification/services"
	log "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/observability/adapters"
	userRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/repository"
	enums "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/card"
	errorsBusiness "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/card/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/card/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/card/dto/response"
	"gorm.io/gorm"
)

type DeleteCardUseCaseIF interface {
	DeleteCard(
		ctx context.Context,
		request request.DeleteCardRequest,
		enterpriseId string,
		userLanguage string,
	) (response.DeleteCardResponse, error)
}

type DeleteCardUseCase struct {
	cardAdapter         adapters.CardAdapter
	cardWriteRepository cardRepository.CardWriteRepositoryIF
	cardReadRepository  cardRepository.CardReadRepositoryIF
	userReadRepository  userRepository.UserReadRepositoryIF
	notificationService services.NotificationService
	log                 log.Logger
}

func NewDeleteCardUseCase(
	cardAdapter adapters.CardAdapter,
	cardWriteRepository cardRepository.CardWriteRepositoryIF,
	cardReadRepository cardRepository.CardReadRepositoryIF,
	userReadRepository userRepository.UserReadRepositoryIF,
	notificationService services.NotificationService,
	log log.Logger,
) DeleteCardUseCaseIF {
	return &DeleteCardUseCase{
		cardAdapter:         cardAdapter,
		cardWriteRepository: cardWriteRepository,
		cardReadRepository:  cardReadRepository,
		userReadRepository:  userReadRepository,
		notificationService: notificationService,
		log:                 log,
	}
}

func (p *DeleteCardUseCase) DeleteCard(
	ctx context.Context,
	request request.DeleteCardRequest,
	enterpriseID string,
	userLanguage string,
) (response.DeleteCardResponse, error) {
	actionDate := time.Now().UTC().In(time.FixedZone("UTC-6", -6*60*60))

	cardUserEmailProjection, getCardAndUserEmailByUserIDErr := p.cardReadRepository.GetCardAndUserEmailByUserID(
		ctx,
		request.UserID,
		request.CardID,
		enterpriseID,
	)

	if getCardAndUserEmailByUserIDErr != nil {
		if errors.Is(getCardAndUserEmailByUserIDErr, gorm.ErrRecordNotFound) {
			return response.DeleteCardResponse{}, errorsBusiness.NewCardNotFoundError(
				request.CardID,
				getCardAndUserEmailByUserIDErr,
			)
		}

		return response.DeleteCardResponse{}, getCardAndUserEmailByUserIDErr
	}

	deletecardr, deleteCardErr := p.cardAdapter.DeleteCard(
		ctx,
		cardUserEmailProjection.ID,
		cardUserEmailProjection.ExternalCardID,
		request.UserID,
		enterpriseID,
	)

	if deleteCardErr != nil {
		return response.DeleteCardResponse{}, errorsBusiness.NewDeleteCardError(enums.ProvidersDeUna, deleteCardErr)
	}

	deleteCardResponse := response.DeleteCardResponse{}

	if deletecardr.Message == "" {
		deleteCardResponse.Status = "Success"
		deleteCardResponse.Message = "Card successfully deleted."
	}

	deleteCardErr = p.cardWriteRepository.DeleteCard(ctx, request.CardID)

	if deleteCardErr != nil {
		return response.DeleteCardResponse{}, errorsBusiness.NewDeleteCardDBError(deleteCardErr)
	}

	formattedActionDate := actionDate.Format("2006-01-02 15:04")

	notifyCardDeletionErr := p.notificationService.NotifyCardDeletion(
		ctx,
		[]constants.NotificationChannel{constants.EmailChannel},
		&cardUserEmailProjection.UserID,
		&cardUserEmailProjection.Email,
		&formattedActionDate,
		&cardUserEmailProjection.LastFour,
		userLanguage,
	)

	if notifyCardDeletionErr != nil {
		p.log.Warn(ctx, fmt.Sprintf("%v", notifyCardDeletionErr))
	}

	return deleteCardResponse, nil
}
