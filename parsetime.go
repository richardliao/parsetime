package parsetime

import (
	"errors"
	"time"
)

var errParse = errors.New("could not parse time")

// ParseInLocation is like time.ParseInLocation.
//
// The result is the given location.
// In the absence of time zone information,
// ParseInLocation interprets the time as in the given location.
//
// The parameter locOffset is added for performance, and its value is the offset corresponding to loc.
//
// For example.
//
// _, locOffset := time.Now().In(loc).Zone()
func ParseInLocation(s string, loc *time.Location, locOffset int) (time.Time, error) {
	t, err := parse([]byte(s), locOffset)
	if err != nil {
		return time.Time{}, nil
	}

	return t.In(loc), nil
}

// Parse is like time.Parse.
//
// The result is the Local location.
// In the absence of a time zone information,
// Parse interprets the time as in UTC.
func Parse(s string) (time.Time, error) {
	return parse([]byte(s), 0)
}

// ParseBytesInLocation is like time.ParseInLocation but accepting bytes with better performance of about 4 ns.
func ParseBytesInLocation(s []byte, loc *time.Location, locOffset int) (time.Time, error) {
	t, err := parse(s, locOffset)
	if err != nil {
		return time.Time{}, nil
	}

	return t.In(loc), nil
}

// ParseBytes is like time.Parse but accepting bytes with better performance of about 4 ns.
func ParseBytes(s []byte) (time.Time, error) {
	return parse(s, 0)
}

func parse(s []byte, locOffset int) (time.Time, error) {
	sLen := len(s)

	if sLen < 10 || s[4] != '-' || s[7] != '-' {
		return time.Time{}, errParse
	}

	var unix int64
	var a0, a1, a2, a3, a4, a5, a6, a7, a8 int

	a0, a1, a2, a3 = int(s[0]-'0'), int(s[1]-'0'), int(s[2]-'0'), int(s[3]-'0')
	if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 || a2 < 0 || a2 > 9 || a3 < 0 || a3 > 9 {
		return time.Time{}, errParse
	}
	year := a0*1e3 + a1*1e2 + a2*1e1 + a3
	month := atoi2MinMax(s[5:7], 1, 12)
	if year == -1 || month == -1 {
		return time.Time{}, errParse
	}
	day := atoi2MinMax(s[8:10], 1, daysIn(time.Month(month), year))
	if day == -1 {
		return time.Time{}, errParse
	}

	// Days since epoc.
	var daysEpoc uint64
	var leap bool

	if year >= unixEpoc && year < unixEpoc+cacheYears {
		daysEpoc = yearDays[year-unixEpoc] + uint64(daysBefore[month-1]) + uint64(day-1)
		leap = yearLeap[year-unixEpoc]
	} else {
		daysEpoc = daysSinceEpoch(year)
		leap = isLeap(year)
	}

	if leap && month >= 3 {
		daysEpoc++
	}

	if sLen == 10 {
		unix = int64(daysEpoc*secondsPerDay) + (absoluteToInternal + internalToUnix)
		return time.Unix(unix-int64(locOffset), 0), nil
	}

	if sLen < 19 || s[13] != ':' || s[16] != ':' || s[10] != 'T' && s[10] != ' ' {
		return time.Time{}, errParse
	}

	hour := atoi2MinMax(s[11:13], 0, 23)
	min := atoi2MinMax(s[14:16], 0, 59)
	sec := atoi2MinMax(s[17:19], 0, 59)
	if hour == -1 || min == -1 || sec == -1 {
		return time.Time{}, errParse
	}

	var nsec, tzSign, tzH, tzM, tzIdx, tzOffset int

	// nsec
	s = s[19:]
	sLen = len(s)
	tzIdx = 0
	if sLen > 1 {
		// .123+08:00
		if s[0] == '.' || s[0] == ',' {
			// Try fast path.
			switch {
			case sLen > 4 && (s[4] == '+' || s[4] == '-' || s[4] == 'z' || s[4] == 'Z'):
				a0, a1, a2 = int(s[1]-'0'), int(s[2]-'0'), int(s[3]-'0')
				if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 || a2 < 0 || a2 > 9 {
					nsec = -1
				} else {
					nsec = (a0*1e2 + a1*1e1 + a2) * 1e6
					tzIdx = 4
				}
			case sLen > 7 && (s[7] == '+' || s[7] == '-' || s[7] == 'z' || s[7] == 'Z'):
				a0, a1, a2, a3, a4 = int(s[1]-'0'), int(s[2]-'0'), int(s[3]-'0'), int(s[4]-'0'), int(s[5]-'0')
				a5 = int(s[6] - '0')
				if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 || a2 < 0 || a2 > 9 || a3 < 0 || a3 > 9 || a4 < 0 || a4 > 9 || a5 < 0 || a5 > 9 {
					nsec = -1
				} else {
					nsec = (a0*1e5 + a1*1e4 + a2*1e3 + a3*1e2 + a4*1e1 + a5) * 1e3
					tzIdx = 7
				}
			case sLen > 10 && (s[10] == '+' || s[10] == '-' || s[10] == 'z' || s[10] == 'Z'):
				a0, a1, a2, a3, a4 = int(s[1]-'0'), int(s[2]-'0'), int(s[3]-'0'), int(s[4]-'0'), int(s[5]-'0')
				a5, a6, a7, a8 = int(s[6]-'0'), int(s[7]-'0'), int(s[8]-'0'), int(s[9]-'0')
				if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 || a2 < 0 || a2 > 9 || a3 < 0 || a3 > 9 || a4 < 0 || a4 > 9 || a5 < 0 || a5 > 9 || a6 < 0 || a6 > 9 || a7 < 0 || a7 > 9 || a8 < 0 || a8 > 9 {
					nsec = -1
				} else {
					nsec = a0*1e8 + a1*1e7 + a2*1e6 + a3*1e5 + a4*1e4 + a5*1e3 + a6*1e2 + a7*1e1 + a8
					tzIdx = 10
				}
			case sLen > 2 && (s[2] == '+' || s[2] == '-' || s[2] == 'z' || s[2] == 'Z'):
				a0 = int(s[1] - '0')
				if a0 < 0 || a0 > 9 {
					nsec = -1
				} else {
					nsec = (a0) * 1e8
					tzIdx = 2
				}
			case sLen > 3 && (s[3] == '+' || s[3] == '-' || s[3] == 'z' || s[3] == 'Z'):
				a0, a1 = int(s[1]-'0'), int(s[2]-'0')
				if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 {
					nsec = -1
				} else {
					nsec = (a0*1e1 + a1) * 1e7
					tzIdx = 3
				}
			case sLen > 5 && (s[5] == '+' || s[5] == '-' || s[5] == 'z' || s[5] == 'Z'):
				a0, a1, a2, a3 = int(s[1]-'0'), int(s[2]-'0'), int(s[3]-'0'), int(s[4]-'0')
				if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 || a2 < 0 || a2 > 9 || a3 < 0 || a3 > 9 {
					nsec = -1
				} else {
					nsec = (a0*1e3 + a1*1e2 + a2*1e1 + a3) * 1e5
					tzIdx = 5
				}
			case sLen > 6 && (s[6] == '+' || s[6] == '-' || s[6] == 'z' || s[6] == 'Z'):
				a0, a1, a2, a3, a4 = int(s[1]-'0'), int(s[2]-'0'), int(s[3]-'0'), int(s[4]-'0'), int(s[5]-'0')
				if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 || a2 < 0 || a2 > 9 || a3 < 0 || a3 > 9 || a4 < 0 || a4 > 9 {
					nsec = -1
				} else {
					nsec = (a0*1e4 + a1*1e3 + a2*1e2 + a3*1e1 + a4) * 1e4
					tzIdx = 6
				}
			case sLen > 8 && (s[8] == '+' || s[8] == '-' || s[8] == 'z' || s[8] == 'Z'):
				a0, a1, a2, a3, a4 = int(s[1]-'0'), int(s[2]-'0'), int(s[3]-'0'), int(s[4]-'0'), int(s[5]-'0')
				a5, a6 = int(s[6]-'0'), int(s[7]-'0')
				if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 || a2 < 0 || a2 > 9 || a3 < 0 || a3 > 9 || a4 < 0 || a4 > 9 || a5 < 0 || a5 > 9 || a6 < 0 || a6 > 9 {
					nsec = -1
				} else {
					nsec = (a0*1e6 + a1*1e5 + a2*1e4 + a3*1e3 + a4*1e2 + a5*1e1 + a6) * 1e2
					tzIdx = 8
				}
			case sLen > 9 && (s[9] == '+' || s[9] == '-' || s[9] == 'z' || s[9] == 'Z'):
				a0, a1, a2, a3, a4 = int(s[1]-'0'), int(s[2]-'0'), int(s[3]-'0'), int(s[4]-'0'), int(s[5]-'0')
				a5, a6, a7 = int(s[6]-'0'), int(s[7]-'0'), int(s[8]-'0')
				if a0 < 0 || a0 > 9 || a1 < 0 || a1 > 9 || a2 < 0 || a2 > 9 || a3 < 0 || a3 > 9 || a4 < 0 || a4 > 9 || a5 < 0 || a5 > 9 || a6 < 0 || a6 > 9 || a7 < 0 || a7 > 9 {
					nsec = -1
				} else {
					nsec = (a0*1e7 + a1*1e6 + a2*1e5 + a3*1e4 + a4*1e3 + a5*1e2 + a6*1e1 + a7) * 1e1
					tzIdx = 9
				}
			default:
				nsec = -1
			}

			// Fallback
			if nsec < 0 && '0' <= s[1] && s[1] <= '9' {
				var val int
				var c byte
				var mult int = 1e9
				for tzIdx = 1; tzIdx < sLen; tzIdx++ {
					if tzIdx > 10 {
						return time.Time{}, errParse
					}

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
		} else if s[0] != 'z' && s[0] != 'Z' && s[0] != '+' && s[0] != '-' {
			return time.Time{}, errParse
		}
	}

	if sLen == 0 || sLen == tzIdx {
		// No tz information.
		unix = int64(daysEpoc*secondsPerDay+uint64(hour*secondsPerHour+min*secondsPerMinute+sec)) + (absoluteToInternal + internalToUnix)
		return time.Unix(unix-int64(locOffset), int64(nsec)), nil
	}

	// Timezone sign.
	switch {
	case sLen == 0 || sLen == tzIdx:
		tzSign = 1
	case s[sLen-1] == 'z' || s[sLen-1] == 'Z' || s[tzIdx] == '+':
		tzSign = 1
	case s[tzIdx] == '-':
		tzSign = -1
	default:
		tzSign = 0
	}

	// Timezone hour and minute.
	s = s[tzIdx:]
	if len(s) > 0 {
		c := s[0]
		if c == 'z' || c == 'Z' {
			tzOffset = 0
		} else {
			if c != '+' && c != '-' {
				return time.Time{}, errParse
			}

			switch len(s) {
			case 6:
				tzH = atoi2MinMax(s[1:3], 0, 14)
				tzM = atoi2MinMax(s[4:6], 0, 59)
				if s[3] != ':' {
					tzH = -1
					tzM = -1
				}
			case 5:
				tzH = atoi2MinMax(s[1:3], 0, 14)
				tzM = atoi2MinMax(s[3:5], 0, 59)
			case 3:
				tzH = atoi2MinMax(s[1:3], 0, 14)
				tzM = 0
			default:
				tzH = -1
				tzM = -1
			}
		}
	}

	if nsec < 0 || tzH == -1 || tzM == -1 || tzSign == 0 {
		return time.Time{}, errParse
	}

	if tzSign == -1 && tzH > 12 {
		return time.Time{}, errParse
	}

	tzOffset = tzSign * (tzH*3600 + tzM*60)

	unix = int64(daysEpoc*secondsPerDay+uint64(hour*secondsPerHour+min*secondsPerMinute+sec)) + (absoluteToInternal + internalToUnix)
	return time.Unix(unix-int64(tzOffset), int64(nsec)), nil
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

// The following code is from the stdlib time.

const (
	secondsPerMinute = 60
	secondsPerHour   = 60 * secondsPerMinute
	secondsPerDay    = 24 * secondsPerHour
	secondsPerWeek   = 7 * secondsPerDay
	daysPer400Years  = 365*400 + 97
	daysPer100Years  = 365*100 + 24
	daysPer4Years    = 365*4 + 1
)

const (
	// The unsigned zero year for internal calculations.
	// Must be 1 mod 400, and times before it will not compute correctly,
	// but otherwise can be changed at will.
	absoluteZeroYear = -292277022399

	// The year of the zero Time.
	// Assumed by the unixToInternal computation below.
	internalYear = 1

	// Offsets to convert between internal and absolute or toUnixUTC times.
	absoluteToInternal int64 = (absoluteZeroYear - internalYear) * 365.2425 * secondsPerDay
	internalToAbsolute       = -absoluteToInternal

	unixToInternal int64 = (1969*365 + 1969/4 - 1969/100 + 1969/400) * secondsPerDay
	internalToUnix int64 = -unixToInternal

	wallToInternal int64 = (1884*365 + 1884/4 - 1884/100 + 1884/400) * secondsPerDay
)

func daysSinceEpoch(year int) uint64 {
	y := uint64(int64(year) - absoluteZeroYear)

	// Add in days from 400-year cycles.
	n := y / 400
	y -= 400 * n
	d := daysPer400Years * n

	// Add in 100-year cycles.
	n = y / 100
	y -= 100 * n
	d += daysPer100Years * n

	// Add in 4-year cycles.
	n = y / 4
	y -= 4 * n
	d += daysPer4Years * n

	// Add in non-leap years.
	n = y
	d += 365 * n

	return d
}

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

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func daysIn(m time.Month, year int) int {
	if m == time.February && isLeap(year) {
		return 29
	}
	return int(daysBefore[m] - daysBefore[m-1])
}

// Cache between 1970 and 2069.
const (
	unixEpoc   = 1970
	cacheYears = 100
)

var (
	yearDays = [cacheYears]uint64{}
	yearLeap = [cacheYears]bool{}
)

func init() {
	for i := 0; i < cacheYears; i++ {
		yearDays[i] = daysSinceEpoch(unixEpoc + i)
		yearLeap[i] = isLeap(unixEpoc + i)
	}
}
