package parsetime

import (
	"fmt"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	now := time.Now()

	_, offset := now.Zone()
	tzOffset := time.Duration(offset) * time.Second

	year := 2023
	month := time.Month(2)
	day := 28
	hour := 15
	min := 00
	sec := 36
	nsec3 := 123000000
	nsec6 := 123456000
	nsec9 := 123456789
	tzH := offset / 3600
	tzM := offset % 3600

	tests := []struct {
		format string
		value  string
		expect time.Time
		err    bool
	}{
		{time.RFC3339Nano, now.Format(time.RFC3339Nano), now, false},
		{time.DateTime, fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec), time.Date(year, month, day, hour, min, sec, 0, time.UTC), false},
		{"2006-01-02T15:04:05", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", year, month, day, hour, min, sec), time.Date(year, month, day, hour, min, sec, 0, time.UTC), false},
		{time.RFC3339, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ", year, month, day, hour, min, sec), time.Date(year, month, day, hour, min, sec, 0, time.UTC), false},
		{"2006-01-02T15:04:05.999999999", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d", year, month, day, hour, min, sec, nsec9), time.Date(year, month, day, hour, min, sec, nsec9, time.UTC), false},
		{"2006-01-02T15:04:05.999999999", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d,%d", year, month, day, hour, min, sec, nsec9), time.Date(year, month, day, hour, min, sec, nsec9, time.UTC), false},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%dz", year, month, day, hour, min, sec, nsec9), time.Date(year, month, day, hour, min, sec, nsec9, time.UTC), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%dZ", year, month, day, hour, min, sec, nsec9), time.Date(year, month, day, hour, min, sec, nsec9, time.UTC), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d:%02d", year, month, day, hour, min, sec, nsec9, tzH, tzM), time.Date(year, month, day, hour, min, sec, nsec9, time.Local), false},
		{"2006-01-02T15:04:05.999999999Z0700", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d%02d", year, month, day, hour, min, sec, nsec9, tzH, tzM), time.Date(year, month, day, hour, min, sec, nsec9, time.Local), false},

		{"2006-01-02T15:04:05.999999999Z07", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d", year, month, day, hour, min, sec, nsec9, tzH), time.Date(year, month, day, hour, min, sec, nsec9, time.Local), false},
		{"2006-01-02T15:04:05.999999999Z07", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d-%02d", year, month, day, hour, min, sec, nsec9, tzH), time.Date(year, month, day, hour, min, sec, nsec9, time.Local).Add(2 * tzOffset), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d-%02d:%02d", year, month, day, hour, min, sec, nsec9, tzH, tzM), time.Date(year, month, day, hour, min, sec, nsec9, time.Local).Add(2 * tzOffset), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d:%02d", year, month, day, hour, min, sec, nsec3, tzH, tzM), time.Date(year, month, day, hour, min, sec, nsec3, time.Local), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d:%02d", year, month, day, hour, min, sec, nsec6, tzH, tzM), time.Date(year, month, day, hour, min, sec, nsec6, time.Local), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d:%02d", year, month, day, hour, min, sec, 0, tzH, tzM), time.Date(year, month, day, hour, min, sec, 0, time.Local), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d:%02d", year, month, day, hour, min, sec, nsec9, 0, 0), time.Date(year, month, day, hour, min, sec, nsec9, time.UTC), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d:%02d", year, month, day, hour, min, sec, nsec9, tzH, 30), time.Date(year, month, day, hour, min, sec, nsec9, time.Local).Add(-30 * time.Minute), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d+%02d:%02d", year, month, day, hour, min, sec, tzH, tzM), time.Date(year, month, day, hour, min, sec, 0, time.Local), false},
		{time.DateOnly, fmt.Sprintf("%04d-%02d-%02d", year, month, day), time.Date(year, month, day, 0, 0, 0, 0, time.UTC), false},

		{"", fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", 2, month, day, hour, min, sec), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dt%02d:%02d:%02d", year, month, day, hour, min, sec), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", year, 13, day, hour, min, sec), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", year, 0, day, hour, min, sec), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", year, -2, day, hour, min, sec), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", year, month, 32, hour, min, sec), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", year, month, 0, hour, min, sec), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", year, month, -2, hour, min, sec), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", 2001, 2, 29, hour, min, sec), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", 2000, 2, 30, hour, min, sec), now, true},

		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", year, month, day, 24, min, sec), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", year, month, day, hour, 60, sec), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", year, month, day, hour, min, 60), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", year, month, day, hour, min, -2), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d:%02d", year, month, day, hour, min, sec, nsec9, 15, tzM), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d%02d:%02d", year, month, day, hour, min, sec, nsec9, -13, tzM), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d:%02d", year, month, day, hour, min, sec, nsec9, tzH, 60), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:aa", year, month, day, hour, min), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02da%d+%02d:%02d", year, month, day, hour, min, sec, nsec6, tzH, tzM), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%da%02d:%02d", year, month, day, hour, min, sec, nsec6, tzH, tzM), now, true},

		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02da%02d", year, month, day, hour, min, sec, nsec6, tzH, tzM), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02da", year, month, day, hour, min, sec, nsec9, tzH), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.aaa+%02d", year, month, day, hour, min, sec, tzH), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d%03d", year, month, day, hour, min, sec, nsec6, tzH, tzM), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%dz%02d:%02d", year, month, day, hour, min, sec, nsec6, tzH, tzM), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d:%02d", year, month, day, hour, min, sec, nsec9, tzH, -10), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d:%02d", year, month, day, hour, min, sec, nsec9*10, tzH, tzM), now, true},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d:%02d", year, month, day, hour, min, sec, nsec9*10, tzH, tzM), now, true},
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

func TestParseInLocation(t *testing.T) {
	now := time.Now()

	loc, _ := time.LoadLocation("Asia/Tokyo")
	_, locOffset := time.Now().In(loc).Zone()

	_, offset := now.In(loc).Zone()
	tzOffset := time.Duration(offset) * time.Second

	year := 2023
	month := time.Month(2)
	day := 28
	hour := 15
	min := 00
	sec := 36
	nsec3 := 123000000
	nsec6 := 123456000
	nsec9 := 123456789
	tzH := offset / 3600
	tzM := offset % 3600

	tests := []struct {
		format string
		value  string
		expect time.Time
		err    bool
	}{
		{time.RFC3339Nano, now.Format(time.RFC3339Nano), now, false},
		{time.DateTime, fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec), time.Date(year, month, day, hour, min, sec, 0, loc), false},
		{"2006-01-02T15:04:05", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", year, month, day, hour, min, sec), time.Date(year, month, day, hour, min, sec, 0, loc), false},
		{time.RFC3339, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ", year, month, day, hour, min, sec), time.Date(year, month, day, hour, min, sec, 0, time.UTC), false},
		{"2006-01-02T15:04:05.999999999", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d", year, month, day, hour, min, sec, nsec9), time.Date(year, month, day, hour, min, sec, nsec9, loc), false},
		{"2006-01-02T15:04:05.999999999", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d,%d", year, month, day, hour, min, sec, nsec9), time.Date(year, month, day, hour, min, sec, nsec9, loc), false},
		{"", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%dz", year, month, day, hour, min, sec, nsec9), time.Date(year, month, day, hour, min, sec, nsec9, time.UTC), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%dZ", year, month, day, hour, min, sec, nsec9), time.Date(year, month, day, hour, min, sec, nsec9, time.UTC), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d:%02d", year, month, day, hour, min, sec, nsec9, tzH, tzM), time.Date(year, month, day, hour, min, sec, nsec9, loc), false},
		{"2006-01-02T15:04:05.999999999Z0700", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d%02d", year, month, day, hour, min, sec, nsec9, tzH, tzM), time.Date(year, month, day, hour, min, sec, nsec9, loc), false},

		{"2006-01-02T15:04:05.999999999Z07", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d", year, month, day, hour, min, sec, nsec9, tzH), time.Date(year, month, day, hour, min, sec, nsec9, loc), false},
		{"2006-01-02T15:04:05.999999999Z07", fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d-%02d", year, month, day, hour, min, sec, nsec9, tzH), time.Date(year, month, day, hour, min, sec, nsec9, loc).Add(2 * tzOffset), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d-%02d:%02d", year, month, day, hour, min, sec, nsec9, tzH, tzM), time.Date(year, month, day, hour, min, sec, nsec9, loc).Add(2 * tzOffset), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d:%02d", year, month, day, hour, min, sec, nsec3, tzH, tzM), time.Date(year, month, day, hour, min, sec, nsec3, loc), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d:%02d", year, month, day, hour, min, sec, nsec6, tzH, tzM), time.Date(year, month, day, hour, min, sec, nsec6, loc), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d:%02d", year, month, day, hour, min, sec, 0, tzH, tzM), time.Date(year, month, day, hour, min, sec, 0, loc), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d:%02d", year, month, day, hour, min, sec, nsec9, 0, 0), time.Date(year, month, day, hour, min, sec, nsec9, time.UTC), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%d+%02d:%02d", year, month, day, hour, min, sec, nsec9, tzH, 30), time.Date(year, month, day, hour, min, sec, nsec9, loc).Add(-30 * time.Minute), false},
		{time.RFC3339Nano, fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d+%02d:%02d", year, month, day, hour, min, sec, tzH, tzM), time.Date(year, month, day, hour, min, sec, 0, loc), false},
		{time.DateOnly, fmt.Sprintf("%04d-%02d-%02d", year, month, day), time.Date(year, month, day, 0, 0, 0, 0, loc), false},
	}

	for i, tt := range tests {
		got, err := ParseInLocation(tt.value, loc, locOffset)
		if tt.err {
			if err == nil {
				t.Fatalf("case %d: expect error got nil", i)
			}
			continue
		}

		if tt.format != "" {
			stdGot, err := time.ParseInLocation(tt.format, tt.value, loc)
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
