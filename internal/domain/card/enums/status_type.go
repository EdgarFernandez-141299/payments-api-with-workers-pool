package enums

import "strings"

type StatusType string

const (
	Default  StatusType = "default"
	Fault    StatusType = "fault"
	Expired  StatusType = "expired"
	Blocked  StatusType = "blocked"
	Active   StatusType = "active"
	Inactive StatusType = "inactive"
)

func (a StatusType) String() string {
	return strings.ToLower(string(a))
}

func (a StatusType) IsValid() bool {
	switch a {
	case Default, Fault, Expired, Blocked, Active, Inactive:
		return true
	}

	return false
}
