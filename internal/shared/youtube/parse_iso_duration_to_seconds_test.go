package youtube

import "testing"

func TestParseISODurationToSeconds(t *testing.T) {
	tests := []struct {
		name      string
		iso       string
		expected  int
		wantError bool
	}{
		{
			name:     "minutes and seconds",
			iso:      "PT4M5S",
			expected: 245,
		},
		{
			name:     "seconds only",
			iso:      "PT45S",
			expected: 45,
		},
		{
			name:     "minutes only",
			iso:      "PT3M",
			expected: 180,
		},
		{
			name:     "hours minutes seconds",
			iso:      "PT1H2M10S",
			expected: 3730,
		},
		{
			name:     "hours only",
			iso:      "PT2H",
			expected: 7200,
		},
		{
			name:      "invalid format",
			iso:       "4M5S",
			wantError: true,
		},
		{
			name:      "empty string",
			iso:       "",
			wantError: true,
		},
		{
			name:     "live stream duration",
			iso:      "PT0S",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseISODurationToSeconds(tt.iso)

			if tt.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
