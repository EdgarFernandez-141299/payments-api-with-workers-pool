package utils

import (
	"errors"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/response"
)

func GetUserEmail(res response.UserProfileInfoDTO) (string, error) {
	if res.Email != "" {
		return res.Email, nil
	}

	return "", errors.New("user email not found")
}

func GetUserPhone(res response.UserProfileInfoDTO) (response.PhoneInfo, error) {
	if res.PrimaryPhone.Number != "" {
		return res.PrimaryPhone, nil
	}

	if res.SecondaryPhone.Number != "" {
		return res.SecondaryPhone, nil
	}

	return response.PhoneInfo{}, errors.New("user phone not found")
}

func GetUserBillingInformation(response response.UserProfileInfoDTO) (response.AddressInfo, error) {
	return response.Address, nil
}
