package value_objects

import (
	"testing"
)

func TestCardInfo_Validate(t *testing.T) {
	tests := []struct {
		name    string
		card    CardInfo
		wantErr bool
	}{
		{
			name:    "valid inputs",
			card:    CardInfo{CardID: "1234567890123456", CVV: "123"},
			wantErr: false,
		},
		{
			name:    "missing CardID",
			card:    CardInfo{CardID: "", CVV: "123"},
			wantErr: true,
		},
		{
			name:    "invalid CVV length - too short",
			card:    CardInfo{CardID: "1234567890123456", CVV: "12"},
			wantErr: true,
		},
		{
			name:    "invalid CVV length - too long",
			card:    CardInfo{CardID: "1234567890123456", CVV: "1234"},
			wantErr: true,
		},
		{
			name:    "missing both CardID and CVV",
			card:    CardInfo{CardID: "", CVV: ""},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.card.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
