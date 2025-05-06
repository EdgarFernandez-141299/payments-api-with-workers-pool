package projections

import "time"

type NotificationCardExpiringSoonProjection struct {
	UserID         string
	LastFour       string
	Email          string
	ExpirationDate time.Time
	EnterpriseID   string
}

type CardUserEmailProjection struct {
	ID             string
	ExternalCardID string
	UserID         string
	LastFour       string
	Email          string
}
