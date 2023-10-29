package parsetime

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	now := time.Now()

	tests := []struct {
		format string
		value  string
		expect time.Time
		err    bool
	}{
		{time.RFC3339Nano, now.Format(time.RFC3339Nano), now, false},
	}

	for i, tt := range tests {

		got, err := Parse(tt.value)
		if tt.err {
			if err == nil {
				t.Fatalf("case %d: expect error got nil", i)
			}
			continue
		}

		if tt.format != "" {
			stdGot, err := time.Parse(tt.format, tt.value)
			if err != nil {
				t.Fatalf("case %d: stdGot error: %s", i, err)
			}

			if !stdGot.Equal(got) {
				t.Fatalf("case %d: got: %+v, std: %+v", i, got, stdGot)
			}
		}

		if !tt.expect.Equal(got) {
			t.Fatalf("case %d: got: %+v, expect: %+v, value: %s", i, got, tt.expect, tt.value)
		}
	}
}
