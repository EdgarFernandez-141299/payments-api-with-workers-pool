package eventsourcing

import (
	"errors"
	"gitlab.com/clubhub.ai1/go-libraries/eventsourcing"
)

type AggregateRoot struct {
	eventsourcing.AggregateRoot
}

var InvalidAggregateID = errors.New("invalid aggregate id")

func IsAggregateIDValid(id string) bool {
	return len(id) > 0
}

func (a *AggregateRoot) WithID(id string) error {
	if !IsAggregateIDValid(id) {
		return InvalidAggregateID
	}

	return a.AggregateRoot.SetID(id)
}
