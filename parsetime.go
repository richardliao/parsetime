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
	if len(s) < 19 || s[4] != '-' || s[7] != '-' || s[13] != ':' || s[16] != ':' || s[10] != 'T' && s[10] != ' ' {
		return time.Time{}, errParse
	}

	var nsec, tzSign, tzH, tzM, tzIdx, tzOffset int

	sLen := len(s)

	year := atoi4(s[0:4])
	month := atoi2MinMax(s[5:7], 1, 12)
	day := atoi2MinMax(s[8:10], 1, daysIn(time.Month(month), year))
	hour := atoi2MinMax(s[11:13], 0, 23)
	min := atoi2MinMax(s[14:16], 0, 59)
	sec := atoi2MinMax(s[17:19], 0, 59)
	if year == -1 || month == -1 || day == -1 || hour == -1 || min == -1 || sec == -1 {
		return time.Time{}, errParse
	}

	// nsec
	tzIdx = 19
	if sLen > 20 && (s[19] == '.' || s[19] == ',') {
		if '0' <= s[20] && s[20] <= '9' {
			var val int
			var c byte
			var mult int = 1e9
			for tzIdx = 20; tzIdx < sLen; tzIdx++ {
				c = s[tzIdx]
				if c >= '0' && c <= '9' {
					val = val*10 + int(c-'0')
					mult /= 10
				} else {
					break
				}
			}
			nsec = val * mult
		}
	}

	// tzH, tzM
	if tzH == -1 || tzM == -1 {
		if s[tzIdx] != '+' && s[tzIdx] != '-' {
			return time.Time{}, errParse
		}

		tzH = atoi2MinMax(s[tzIdx+1:tzIdx+3], 0, 23)

		tzmIdx := 3
		if s[tzIdx+3] == ':' {
			tzmIdx++
		}

		tzM = atoi2MinMax(s[tzIdx+tzmIdx:tzIdx+tzmIdx+2], 0, 59)
	}

	// tzSign
	switch {
	case s[sLen-6] == '+' || s[sLen-5] == '+':
		tzSign = 1
	case s[sLen-6] == '-' || s[sLen-5] == '-':
		tzSign = -1
	case s[sLen-1] == 'z' || s[sLen-1] == 'Z':
		tzSign = 1
	default:
		tzSign = 0
	}

	if nsec == -1 || tzH == -1 || tzM == -1 || tzSign == 0 {
		return time.Time{}, errParse
	}

	tzOffset = tzSign * (tzH*3600 + tzM*60)

	t := time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.Local)
	t = t.Add(-time.Duration(tzOffset) * time.Second)
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
