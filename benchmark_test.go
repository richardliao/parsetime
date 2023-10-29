package parsetime_test

import (
	"github.com/richardliao/parsetime"
	"testing"
	"time"
)

func BenchmarkRFC3339NanoBytes(b *testing.B) {
	now := []byte(time.Now().Local().Format(time.RFC3339Nano))

	for i := 0; i < b.N; i++ {
		if _, err := parsetime.ParseBytes(now); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMultiFormat(b *testing.B) {
	now := time.Now().Local().Format(time.RFC3339Nano)

	for i := 0; i < b.N; i++ {
		if _, err := parsetime.Parse(now); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRFC3339Nano(b *testing.B) {
	now := time.Now().Local().Format(time.RFC3339Nano)

	for i := 0; i < b.N; i++ {
		if _, err := parsetime.Parse(now); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRFC3339(b *testing.B) {
	now := time.Now().Local().Format(time.RFC3339)

	for i := 0; i < b.N; i++ {
		if _, err := parsetime.Parse(now); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDateTime(b *testing.B) {
	now := time.Now().Local().Format(time.DateTime)

	for i := 0; i < b.N; i++ {
		if _, err := parsetime.Parse(now); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDateOnly(b *testing.B) {
	now := time.Now().Local().Format(time.DateOnly)

	for i := 0; i < b.N; i++ {
		if _, err := parsetime.Parse(now); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNonStandardFormat(b *testing.B) {
	now := time.Now().Local().Format("2006-01-02 15:04:05.999999999Z07:00")

	for i := 0; i < b.N; i++ {
		if _, err := parsetime.Parse(now); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRFC3339InLocation(b *testing.B) {
	now := time.Now().Local().Format(time.RFC3339)
	_, locOffset := time.Now().In(time.Local).Zone()

	for i := 0; i < b.N; i++ {
		if _, err := parsetime.ParseInLocation(now, time.Local, locOffset); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoMultiFormat(b *testing.B) {
	now := time.Now().Local().Format(time.RFC3339Nano)

	for i := 0; i < b.N; i++ {
		if _, err := time.Parse(time.DateOnly, now); err != nil {
			if _, err := time.Parse(time.DateTime, now); err != nil {
				if _, err := time.Parse(time.RFC3339, now); err != nil {
					if _, err := time.Parse(time.RFC3339Nano, now); err != nil {
						b.Fatal(err)
					}
				}
			}
		}
	}
}

func BenchmarkGoRFC3339Nano(b *testing.B) {
	now := time.Now().Local().Format(time.RFC3339Nano)

	for i := 0; i < b.N; i++ {
		if _, err := time.Parse(time.RFC3339Nano, now); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoRFC3339(b *testing.B) {
	now := time.Now().Local().Format(time.RFC3339)

	for i := 0; i < b.N; i++ {
		if _, err := time.Parse(time.RFC3339, now); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoDateTime(b *testing.B) {
	now := time.Now().Local().Format(time.DateTime)

	for i := 0; i < b.N; i++ {
		if _, err := time.Parse(time.DateTime, now); err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkGoDateOnly(b *testing.B) {
	now := time.Now().Local().Format(time.DateOnly)

	for i := 0; i < b.N; i++ {
		if _, err := time.Parse(time.DateOnly, now); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoNonStandardFormat(b *testing.B) {
	now := time.Now().Local().Format("2006-01-02 15:04:05.999999999Z07:00")

	for i := 0; i < b.N; i++ {
		if _, err := time.Parse("2006-01-02 15:04:05.999999999Z07:00", now); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoRFC3339InLocation(b *testing.B) {
	now := time.Now().Local().Format(time.RFC3339)

	for i := 0; i < b.N; i++ {
		if _, err := time.ParseInLocation(time.RFC3339, now, time.Local); err != nil {
			b.Fatal(err)
		}
	}
}
