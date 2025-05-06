package enums

import "strings"

type StatusRoute string

const (
	Pending StatusRoute = "PENDING"
	Active  StatusRoute = "ACTIVE"
)

func (s StatusRoute) String() string {
	return strings.ToUpper(string(s))
}

func (s StatusRoute) IsValid() bool {
	switch s {
	case Pending, Active:
		return true
	}

	return false
}
