package enums

import "testing"

func TestAssociatedOriginString(t *testing.T) {
	tests := []struct {
		input    AssociatedOrigin
		expected string
	}{
		{Downpayment, "DOWNPAYMENT"},
		{Loan, "LOAN"},
		{Club, "CLUB"},
		{Booking, "BOOKING"},
	}

	for _, test := range tests {
		if result := test.input.String(); result != test.expected {
			t.Errorf("expected %s, got %s", test.expected, result)
		}
	}
}

func TestAssociatedOrigin_IsValid(t *testing.T) {
	tests := []struct {
		input    AssociatedOrigin
		expected bool
	}{
		{Downpayment, true},
		{Loan, true},
		{Club, true},
		{Booking, true},
		{"INVALID", false},
	}

	for _, test := range tests {
		if result := test.input.IsValid(); result != test.expected {
			t.Errorf("expected %v, got %v", test.expected, result)
		}
	}
}

func TestNewAssociatedOrigin(t *testing.T) {

	tests := []struct {
		input    string
		expected AssociatedOrigin
		err      bool
	}{
		{"DOWNPAYMENT", Downpayment, false},
		{"LOAN", Loan, false},
		{"CLUB", Club, false},
		{"BOOKING", Booking, false},
		{"INVALID", "", true},
	}

	for _, test := range tests {
		result, err := NewAssociatedOrigin(test.input)
		if (err != nil) != test.err {
			t.Errorf("expected error: %v, got: %v", test.err, err)
		}
		if result != test.expected {
			t.Errorf("expected %v, got %v", test.expected, result)
		}
	}

}
