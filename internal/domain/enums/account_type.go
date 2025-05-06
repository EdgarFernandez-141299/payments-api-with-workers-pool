package enums

import "strings"

type AccountType string

const (
	Concentrators AccountType = "CONCENTRATORS"
	Payers        AccountType = "PAYERS"
	Mixed         AccountType = "MIXED"
)

func (a AccountType) String() string {
	return strings.ToUpper(string(a))
}

func (a AccountType) IsValid() bool {
	switch a {
	case Concentrators, Payers, Mixed:
		return true
	}

	return false
}
