package value_objects

import (
	"errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
)

type AssociatedOrigin struct {
	Type enums.AssociatedOrigin
}

func NewAssociatedOrigin(associatedOriginType enums.AssociatedOrigin) AssociatedOrigin {
	return AssociatedOrigin{
		Type: associatedOriginType,
	}
}

func NewFromAssociatedOriginString(origin string) (AssociatedOrigin, error) {
	associatedOriginTypeEnum, err := enums.NewAssociatedOrigin(origin)

	if err != nil {
		return AssociatedOrigin{}, err
	}

	return NewAssociatedOrigin(associatedOriginTypeEnum), nil
}

func (a AssociatedOrigin) Validate() error {
	if !a.Type.IsValid() {
		return errors.New("type is not valid")
	}

	return nil
}
