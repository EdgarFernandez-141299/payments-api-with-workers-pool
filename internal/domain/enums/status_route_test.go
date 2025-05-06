package enums

import "testing"

func TestStatusRouteString(t *testing.T) {
	tests := []struct {
		name string
		s    StatusRoute
		want string
	}{
		{"Pending", Pending, "PENDING"},
		{"Active", Active, "ACTIVE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("StatusRoute.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusRoute_IsValid(t *testing.T) {
	tests := []struct {
		name string
		s    StatusRoute
		want bool
	}{
		{"ValidPending", Pending, true},
		{"ValidActive", Active, true},
		{"InvalidStatus", StatusRoute("INVALID"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsValid(); got != tt.want {
				t.Errorf("StatusRoute.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
