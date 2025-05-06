package adapters

import (
	"testing"
)

func TestMetadata_Add(t *testing.T) {
	tests := []struct {
		name         string
		initial      Metadata
		key          string
		value        string
		expectResult Metadata
	}{
		{
			name:         "Add to nil Metadata",
			initial:      nil,
			key:          "key1",
			value:        "value1",
			expectResult: Metadata{"key1": "value1"},
		},
		{
			name:         "Add to empty Metadata",
			initial:      Metadata{},
			key:          "key2",
			value:        "value2",
			expectResult: Metadata{"key2": "value2"},
		},
		{
			name:         "Overwrite existing key",
			initial:      Metadata{"key1": "value1"},
			key:          "key1",
			value:        "value2",
			expectResult: Metadata{"key1": "value2"},
		},
		{
			name:         "Add new key to non-empty Metadata",
			initial:      Metadata{"key1": "value1"},
			key:          "key2",
			value:        "value2",
			expectResult: Metadata{"key1": "value1", "key2": "value2"},
		},
		{
			name:         "Add key with empty string value",
			initial:      Metadata{},
			key:          "key1",
			value:        "",
			expectResult: Metadata{"key1": ""},
		},
		{
			name:         "Add key with empty string key and value",
			initial:      Metadata{},
			key:          "",
			value:        "",
			expectResult: Metadata{"": ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.initial
			m.Add(tt.key, tt.value)

			if len(m) != len(tt.expectResult) {
				t.Errorf("expected Metadata of size %d, got %d", len(tt.expectResult), len(m))
			}

			for k, v := range tt.expectResult {
				if m[k] != v {
					t.Errorf("expected key %q to have value %q, got %q", k, v, m[k])
				}
			}
		})
	}
}
