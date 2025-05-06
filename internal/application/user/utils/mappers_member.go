package utils

import (
	"errors"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/response"
)

func GetMainEmailMember(emails []response.EmailDTO) (string, error) {
	for _, e := range emails {
		if e.IsDefault {
			return e.Email, nil
		}
	}

	return "", errors.New("main email not found")
}

func GetMainPhoneMember(phones []response.PhoneDTO) (string, error) {
	for _, p := range phones {
		if p.IsDefault {
			return p.Number, nil
		}
	}

	return "", errors.New("main phone not found")
}

func GetBillingInformationMember(billingInformation response.BillingInformationDTO) (response.BillingInformationDTO, error) {

	if billingInformation != (response.BillingInformationDTO{}) {
		return billingInformation, nil
	}

	return response.BillingInformationDTO{}, nil
}
