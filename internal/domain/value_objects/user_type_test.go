package value_objects

import (
	"testing"
)

func TestNewUserType(t *testing.T) {
	type testCase struct {
		name          string
		input         UserType
		expectedValue UserType
	}

	tests := []testCase{
		{
			name:          "ValidMember",
			input:         Member,
			expectedValue: Member,
		},
		{
			name:          "ValidEmployee",
			input:         Employee,
			expectedValue: Employee,
		},
		{
			name:          "InvalidType",
			input:         Invalid,
			expectedValue: Invalid,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := NewUserType(tc.input)
			if !result.Equals(tc.expectedValue) {
				t.Errorf("expected %v, got %v", tc.expectedValue, result)
			}
		})
	}
}

func TestNewUserTypeFromString(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expectedValue UserType
	}

	tests := []testCase{
		{
			name:          "ValidMemberString",
			input:         "member",
			expectedValue: Member,
		},
		{
			name:          "ValidEmployeeString",
			input:         "employee",
			expectedValue: Employee,
		},
		{
			name:          "InvalidString",
			input:         "unknown",
			expectedValue: Invalid,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := NewUserTypeFromString(tc.input)
			if !result.Equals(tc.expectedValue) {
				t.Errorf("expected %v, got %v", tc.expectedValue, result)
			}
		})
	}
}
