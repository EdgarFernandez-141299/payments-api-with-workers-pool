package enums

import "testing"

func TestAccountTypeString(t *testing.T) {
	tests := []struct {
		name string
		a    AccountType
		want string
	}{
		{"Concentrators", Concentrators, "CONCENTRATORS"},
		{"Payers", Payers, "PAYERS"},
		{"Mixed", Mixed, "MIXED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("AccountType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		a    AccountType
		want bool
	}{
		{"Valid Concentrators", Concentrators, true},
		{"Valid Payers", Payers, true},
		{"Valid Mixed", Mixed, true},
		{"Invalid AccountType", AccountType("INVALID"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.IsValid(); got != tt.want {
				t.Errorf("AccountType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
