package enums

import (
	"testing"
)

func TestStatusType_String(t *testing.T) {
	tests := []struct {
		input    StatusType
		expected string
	}{
		{Default, "default"},
		{Fault, "fault"},
		{Expired, "expired"},
		{Blocked, "blocked"},
		{Active, "active"},
		{Inactive, "inactive"},
	}

	for _, test := range tests {
		result := test.input.String()
		if result != test.expected {
			t.Errorf("Expected %s, but got %s", test.expected, result)
		}
	}
}

func TestStatusType_IsValid(t *testing.T) {
	tests := []struct {
		input    StatusType
		expected bool
	}{
		{Default, true},
		{Fault, true},
		{Expired, true},
		{Blocked, true},
		{Active, true},
		{Inactive, true},
		{"INVALID", false},
	}

	for _, test := range tests {
		result := test.input.IsValid()
		if result != test.expected {
			t.Errorf("Expected %v, but got %v", test.expected, result)
		}
	}
}
