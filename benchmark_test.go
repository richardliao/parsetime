package parsetime_test

import (
	"github.com/richardliao/parsetime"
	"testing"
	"time"
)

func BenchmarkRFC3339Nano(b *testing.B) {
	now := time.Now().Local().Format(time.RFC3339Nano)

	for i := 0; i < b.N; i++ {
		if _, err := parsetime.Parse(now); err != nil {
			b.Fatal(err)
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
