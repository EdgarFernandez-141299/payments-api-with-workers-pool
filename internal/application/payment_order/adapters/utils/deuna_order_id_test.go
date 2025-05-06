package utils

import (
	"testing"
)

func TestExtractFromDeunaOrderID(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		expected  DeunaOrderID
	}{
		{
			name:      "valid input",
			input:     "ref123-pay456",
			wantError: false,
			expected: DeunaOrderID{
				composedOrderID:  "ref123-pay456",
				referenceOrderID: "ref123",
				paymentID:        "pay456",
			},
		},
		{
			name:      "missing paymentID",
			input:     "ref123-",
			wantError: false,
			expected: DeunaOrderID{
				composedOrderID:  "ref123-",
				referenceOrderID: "ref123",
				paymentID:        "",
			},
		},
		{
			name:      "missing referenceOrderID",
			input:     "-pay456",
			wantError: false,
			expected: DeunaOrderID{
				composedOrderID:  "-pay456",
				referenceOrderID: "",
				paymentID:        "pay456",
			},
		},
		{
			name:      "no separator",
			input:     "ref123pay456",
			wantError: true,
		},
		{
			name:      "extra separators",
			input:     "ref123-pay456-extra",
			wantError: false,
			expected: DeunaOrderID{
				composedOrderID:  "ref123-pay456-extra",
				referenceOrderID: "ref123-pay456",
				paymentID:        "extra",
			},
		},
		{
			name:      "empty input",
			input:     "",
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ExtractFromDeunaOrderID(tc.input)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error but got none for input: '%s'", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input '%s': %v", tc.input, err)
				}
				if got.composedOrderID != tc.expected.composedOrderID {
					t.Errorf("expected composedOrderID '%s', got '%s'", tc.expected.composedOrderID, got.composedOrderID)
				}
				if got.referenceOrderID != tc.expected.referenceOrderID {
					t.Errorf("expected referenceOrderID '%s', got '%s'", tc.expected.referenceOrderID, got.referenceOrderID)
				}
				if got.paymentID != tc.expected.paymentID {
					t.Errorf("expected paymentID '%s', got '%s'", tc.expected.paymentID, got.paymentID)
				}
			}
		})
	}
}

func TestNewDeunaOrderID(t *testing.T) {
	tests := []struct {
		name      string
		orderID   string
		paymentID string
		expected  DeunaOrderID
	}{
		{
			name:      "valid input",
			orderID:   "ref123",
			paymentID: "pay456",
			expected: DeunaOrderID{
				composedOrderID:  "ref123-pay456",
				referenceOrderID: "ref123",
				paymentID:        "pay456",
			},
		},
		{
			name:      "empty orderID",
			orderID:   "",
			paymentID: "pay456",
			expected: DeunaOrderID{
				composedOrderID:  "-pay456",
				referenceOrderID: "",
				paymentID:        "pay456",
			},
		},
		{
			name:      "empty paymentID",
			orderID:   "ref123",
			paymentID: "",
			expected: DeunaOrderID{
				composedOrderID:  "ref123-",
				referenceOrderID: "ref123",
				paymentID:        "",
			},
		},
		{
			name:      "empty orderID and paymentID",
			orderID:   "",
			paymentID: "",
			expected: DeunaOrderID{
				composedOrderID:  "-",
				referenceOrderID: "",
				paymentID:        "",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := NewDeunaOrderID(tc.orderID, tc.paymentID)
			if got.composedOrderID != tc.expected.composedOrderID {
				t.Errorf("expected composedOrderID '%s', got '%s'", tc.expected.composedOrderID, got.composedOrderID)
			}
			if got.referenceOrderID != tc.expected.referenceOrderID {
				t.Errorf("expected referenceOrderID '%s', got '%s'", tc.expected.referenceOrderID, got.referenceOrderID)
			}
			if got.paymentID != tc.expected.paymentID {
				t.Errorf("expected paymentID '%s', got '%s'", tc.expected.paymentID, got.paymentID)
			}

		})
	}
}
