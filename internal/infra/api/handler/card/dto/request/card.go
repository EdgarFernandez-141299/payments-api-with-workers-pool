package request

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/notification/constants"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/card/enums"
)

type CardRequest struct {
	UserID         string `json:"user_id" validate:"required,lte=36"`
	CardId         string `json:"card_token_id" validate:"required,lte=50"`
	Alias          string `json:"alias" validate:"omitempty,lte=50"`
	Status         string `json:"status" validate:"required,lte=30,status"`
	CardType       string `json:"card_type" validate:"required"`
	CardBrand      string `json:"company" validate:"required"`
	LastFour       string `json:"last_four" validate:"required"`
	FirstSix       string `json:"first_six" validate:"required"`
	ExpirationDate string `json:"expiration_date" validate:"required"`
	IsDefault      bool   `json:"is_default"`
	IsRecurrent    bool   `json:"is_recurrent"`
}

type DeleteCardRequest struct {
	CardID string `json:"card_id" validate:"required,lte=50"`
	UserID string `json:"user_id" validate:"required,lte=50"`
}

type NotificationCardExpiringSoonRequestDTO struct {
	NotificationChannels []constants.NotificationChannel `json:"notification_channels" validate:"required"`
}

var validate *validator.Validate

func statusValidation(fl validator.FieldLevel) bool {
	status := strings.TrimSpace(fl.Field().String())
	return strings.ToLower(status) == strings.ToLower(enums.Active.String()) || status == ""
}

func (c *CardRequest) Validate() error {
	validate = validator.New()

	validate.RegisterValidation("status", statusValidation)

	err := validate.Struct(c)

	if err != nil {
		return err
	}

	c.Status = strings.TrimSpace(c.Status)
	if c.Status == "" {
		c.Status = strings.ToLower(enums.Active.String())
	}

	return nil
}

func (c *DeleteCardRequest) Validate() error {
	err := validator.New().Struct(c)

	if err != nil {
		return err
	}

	return nil
}

func (ncesrdto *NotificationCardExpiringSoonRequestDTO) Validate() error {
	err := validator.New().Struct(ncesrdto)

	if err != nil || len(ncesrdto.NotificationChannels) == 0 {
		return errors.New("notification_channels is required")
	}

	for _, channel := range ncesrdto.NotificationChannels {
		if !constants.AllowedNotificationChannels[channel] {
			return fmt.Errorf("invalid notification channel: %s", channel)
		}
	}

	return nil
}
