package parsetime

import (
	"errors"
	"time"
)

var errParse = errors.New("could not parse time")

// Parse is like time.Parse.
//
// The result is the Local location.
// In the absence of a time zone information,
// Parse interprets the time as in UTC.
func Parse(s string) (time.Time, error) {
	return parse([]byte(s))
}

func parse(s []byte) (time.Time, error) {
	ok := true

	// Parse the date and time.
	if len(s) < len("2006-01-02T15:04:05") {
		return time.Time{}, errParse
	}
	year := atoi4(s[0:4])                                           // e.g., 2006
	month := atoi2MinMax(s[5:7], 1, 12)                             // e.g., 01
	day := atoi2MinMax(s[8:10], 1, daysIn(time.Month(month), year)) // e.g., 02
	hour := atoi2MinMax(s[11:13], 0, 23)                            // e.g., 15
	min := atoi2MinMax(s[14:16], 0, 59)                             // e.g., 04
	sec := atoi2MinMax(s[17:19], 0, 59)                             // e.g., 05
	if !ok || !(s[4] == '-' && s[7] == '-' && s[10] == 'T' && s[13] == ':' && s[16] == ':') {
		return time.Time{}, errParse
	}
	s = s[19:]

	// Parse the fractional second.
	var nsec int
	if len(s) >= 2 && s[0] == '.' && isDigit(s, 1) {
		n := 2
		for ; n < len(s) && isDigit(s, n); n++ {
		}
		nsec, _, _ = parseNanoseconds(s, n)
		s = s[n:]
	}

	// Parse the time zone.
	t := time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.UTC)
	if len(s) != 1 || s[0] != 'Z' {
		if len(s) != len("-07:00") {
			return time.Time{}, errParse
		}
		hr := atoi2MinMax(s[1:3], 0, 23) // e.g., 07
		mm := atoi2MinMax(s[4:6], 0, 59) // e.g., 00
		if !ok || !((s[0] == '-' || s[0] == '+') && s[3] == ':') {
			return time.Time{}, errParse
		}
		zoneOffset := (hr*60 + mm) * 60
		if s[0] == '-' {
			zoneOffset *= -1
		}

		t = t.Add(-time.Duration(zoneOffset) * time.Second)
	}
	return t, nil
}

func atoi2MinMax(s []byte, min, max int) (x int) {
	a0, a1 := int(s[0]-'0'), int(s[1]-'0')
	if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 {
		return -1
	}
	x = a0*1e1 + a1
	if x < min || max < x {
		return -1
	}
	return x
}

func atoi4(s []byte) (x int) {
	a0, a1, a2, a3 := int(s[0]-'0'), int(s[1]-'0'), int(s[2]-'0'), int(s[3]-'0')
	if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 || a2 < 0 || a2 > 9 || a3 < 0 || a3 > 9 {
		return -1
	}
	x = a0*1e3 + a1*1e2 + a2*1e1 + a3
	return x
}

// The following code is almost from the stdlib time.

func isDigit(s []byte, i int) bool {
	if len(s) <= i {
		return false
	}
	c := s[i]
	return '0' <= c && c <= '9'
}

func parseNanoseconds(value []byte, nbytes int) (ns int, rangeErrString string, err error) {
	if !commaOrPeriod(value[0]) {
		err = errParse
		return
	}
	if nbytes > 10 {
		value = value[:10]
		nbytes = 10
	}
	if ns, err = atoi(value[1:nbytes]); err != nil {
		return
	}
	if ns < 0 {
		rangeErrString = "fractional second"
		return
	}
	// We need nanoseconds, which means scaling by the number
	// of missing digits in the format, maximum length 10.
	scaleDigits := 10 - nbytes
	for i := 0; i < scaleDigits; i++ {
		ns *= 10
	}
	return
}

func commaOrPeriod(b byte) bool {
	return b == '.' || b == ','
}

func atoi(s []byte) (x int, err error) {
	neg := false
	if len(s) > 0 && (s[0] == '-' || s[0] == '+') {
		neg = s[0] == '-'
		s = s[1:]
	}
	q, rem, err := leadingInt(s)
	x = int(q)
	if err != nil || len(rem) > 0 {
		return 0, errParse
	}
	if neg {
		x = -x
	}
	return x, nil
}

func leadingInt(s []byte) (x uint64, rem []byte, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > 1<<63/10 {
			// overflow
			return 0, rem, errParse
		}
		x = x*10 + uint64(c) - '0'
		if x > 1<<63 {
			// overflow
			return 0, rem, errParse
		}
	}
	return x, s[i:], nil
}

func daysIn(m time.Month, year int) int {
	if m == time.February && isLeap(year) {
		return 29
	}
	return int(daysBefore[m] - daysBefore[m-1])
}

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// daysBefore[m] counts the number of days in a non-leap year
// before month m begins. There is an entry for m=12, counting
// the number of days before January of next year (365).
var daysBefore = [...]int32{
	0,
	31,
	31 + 28,
	31 + 28 + 31,
	31 + 28 + 31 + 30,
	31 + 28 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + 31,
}
