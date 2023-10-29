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
		// fast path
		switch {
		case sLen > 22 && (s[23] == '+' || s[23] == '-' || s[23] == 'z' || s[23] == 'Z'):
			nsec = atoi3(s[20:23]) * 1e6
			tzIdx = 23
		case sLen > 25 && (s[26] == '+' || s[26] == '-' || s[26] == 'z' || s[26] == 'Z'):
			nsec = atoi6(s[20:26]) * 1e3
			tzIdx = 26
		case sLen > 28 && (s[29] == '+' || s[29] == '-' || s[29] == 'z' || s[29] == 'Z'):
			nsec = atoi9(s[20:29])
			tzIdx = 29
		case sLen > 20 && (s[21] == '+' || s[21] == '-' || s[21] == 'z' || s[21] == 'Z'):
			nsec = atoi1(s[20:21]) * 1e8
			tzIdx = 21
		case sLen > 21 && (s[22] == '+' || s[22] == '-' || s[22] == 'z' || s[22] == 'Z'):
			nsec = atoi2(s[20:22]) * 1e6
			tzIdx = 22
		case sLen > 23 && (s[24] == '+' || s[24] == '-' || s[24] == 'z' || s[24] == 'Z'):
			nsec = atoi4(s[20:24]) * 1e5
			tzIdx = 24
		case sLen > 24 && (s[25] == '+' || s[25] == '-' || s[25] == 'z' || s[25] == 'Z'):
			nsec = atoi5(s[20:25]) * 1e4
			tzIdx = 25
		case sLen > 26 && (s[27] == '+' || s[27] == '-' || s[27] == 'z' || s[27] == 'Z'):
			nsec = atoi7(s[20:27]) * 1e2
			tzIdx = 27
		case sLen > 27 && (s[28] == '+' || s[28] == '-' || s[28] == 'z' || s[28] == 'Z'):
			nsec = atoi8(s[20:28]) * 1e1
			tzIdx = 28
		default:
			nsec = -1
		}

		// fallback
		if nsec == -1 && '0' <= s[20] && s[20] <= '9' {
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
	switch {
	case s[sLen-3] == ':':
		if s[sLen-6] == '+' || s[sLen-6] == '-' {
			tzH = atoi2MinMax(s[sLen-5:sLen-3], 0, 23)
			tzM = atoi2MinMax(s[sLen-2:sLen], 0, 59)
		}
	case s[sLen-4] == '+' || s[sLen-4] == '-':
		tzH = atoi2MinMax(s[sLen-4:sLen-2], 0, 23)
		tzM = atoi2MinMax(s[sLen-2:sLen], 0, 59)
	case s[sLen-2] == '+' || s[sLen-2] == '-':
		tzH = atoi2MinMax(s[sLen-2:sLen], 0, 23)
		tzM = 0
	case s[sLen-1] == 'z' || s[sLen-1] == 'Z':
		tzH = 0
		tzM = 0
	default:
		tzH = -1
		tzM = -1
	}

	// fallback
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

	t := time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.UTC)
	t = t.Add(-time.Duration(tzOffset) * time.Second)
	return t, nil
}

func atoi2MinMax(s []byte, min, max int) (x int) {
	_ = s[1]
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

func atoi1(s []byte) (x int) {
	_ = s[0]
	a0 := int(s[0] - '0')
	if a0 < 0 || a0 > 9 {
		return -1
	}
	return a0 * 1
}

func atoi2(s []byte) (x int) {
	_ = s[1]
	a0, a1 := int(s[0]-'0'), int(s[1]-'0')
	if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 {
		return -1
	}
	return a0*1e1 + a1*1
}

func atoi3(s []byte) (x int) {
	_ = s[2]
	a0, a1, a2 := int(s[0]-'0'), int(s[1]-'0'), int(s[2]-'0')
	if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 || a2 < 0 || a2 > 9 {
		return -1
	}
	return a0*1e2 + a1*1e1 + a2*1
}

func atoi4(s []byte) (x int) {
	_ = s[3]
	a0, a1, a2, a3 := int(s[0]-'0'), int(s[1]-'0'), int(s[2]-'0'), int(s[3]-'0')
	if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 || a2 < 0 || a2 > 9 || a3 < 0 || a3 > 9 {
		return -1
	}
	return a0*1e3 + a1*1e2 + a2*1e1 + a3*1
}

func atoi5(s []byte) (x int) {
	_ = s[4]
	a0, a1, a2, a3, a4 := int(s[0]-'0'), int(s[1]-'0'), int(s[2]-'0'), int(s[3]-'0'), int(s[4]-'0')
	if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 || a2 < 0 || a2 > 9 || a3 < 0 || a3 > 9 || a4 < 0 || a4 > 9 {
		return -1
	}

	return a0*1e4 + a1*1e3 + a2*1e2 + a3*1e1 + a4*1
}

func atoi6(s []byte) (x int) {
	_ = s[5]
	a0, a1, a2, a3, a4 := int(s[0]-'0'), int(s[1]-'0'), int(s[2]-'0'), int(s[3]-'0'), int(s[4]-'0')
	if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 || a2 < 0 || a2 > 9 || a3 < 0 || a3 > 9 || a4 < 0 || a4 > 9 {
		return -1
	}
	a5 := int(s[5] - '0')
	if a5 < 0 || a5 > 9 {
		return -1
	}

	return a0*1e5 + a1*1e4 + a2*1e3 + a3*1e2 + a4*1e1 + a5*1
}

func atoi7(s []byte) (x int) {
	_ = s[6]
	a0, a1, a2, a3, a4 := int(s[0]-'0'), int(s[1]-'0'), int(s[2]-'0'), int(s[3]-'0'), int(s[4]-'0')
	if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 || a2 < 0 || a2 > 9 || a3 < 0 || a3 > 9 || a4 < 0 || a4 > 9 {
		return -1
	}
	a5, a6 := int(s[5]-'0'), int(s[6]-'0')
	if a5 < 0 || a5 > 9 || a6 < 0 || a6 > 9 {
		return -1
	}

	return a0*1e6 + a1*1e5 + a2*1e4 + a3*1e3 + a4*1e2 + a5*1e1 + a6*1
}

func atoi8(s []byte) (x int) {
	_ = s[7]
	a0, a1, a2, a3, a4 := int(s[0]-'0'), int(s[1]-'0'), int(s[2]-'0'), int(s[3]-'0'), int(s[4]-'0')
	if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 || a2 < 0 || a2 > 9 || a3 < 0 || a3 > 9 || a4 < 0 || a4 > 9 {
		return -1
	}
	a5, a6, a7 := int(s[5]-'0'), int(s[6]-'0'), int(s[7]-'0')
	if a5 < 0 || a5 > 9 || a6 < 0 || a6 > 9 || a7 < 0 || a7 > 9 {
		return -1
	}

	return a0*1e7 + a1*1e6 + a2*1e5 + a3*1e4 + a4*1e3 + a5*1e2 + a6*1e1 + a7*1
}

func atoi9(s []byte) (x int) {
	_ = s[8]
	a0, a1, a2, a3, a4 := int(s[0]-'0'), int(s[1]-'0'), int(s[2]-'0'), int(s[3]-'0'), int(s[4]-'0')
	if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 || a2 < 0 || a2 > 9 || a3 < 0 || a3 > 9 || a4 < 0 || a4 > 9 {
		return -1
	}
	a5, a6, a7, a8 := int(s[5]-'0'), int(s[6]-'0'), int(s[7]-'0'), int(s[8]-'0')
	if a5 < 0 || a5 > 9 || a6 < 0 || a6 > 9 || a7 < 0 || a7 > 9 || a8 < 0 || a8 > 9 {
		return -1
	}

	return a0*1e8 + a1*1e7 + a2*1e6 + a3*1e5 + a4*1e4 + a5*1e3 + a6*1e2 + a7*1e1 + a8*1
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
