# Overview

ParseTime is a Go package dedicated to quickly parse RFC3339 or similar time strings.

The time.Parse in Go stdlib is fast enough in most cases, but not when processing massive time series data.
The time string generated in many languages is not the standard RFC3339, but a variant of it.
When it is not a standard RFC3339 string (e.g. 2006-01-02T15:04:05.999), 
or multiple uncertain time formats need to be tried, 
the resources consumed by the parsing time are very considerable.

```
// parsetime
BenchmarkRFC3339NanoBytes    	257910897	        23.21 ns/op
BenchmarkMultiFormat         	219147885	        27.09 ns/op
BenchmarkRFC3339Nano         	222247737	        27.01 ns/op
BenchmarkRFC3339             	253692507	        23.56 ns/op
BenchmarkDateTime            	346593026	        17.33 ns/op
BenchmarkDateOnly            	470903756	        12.57 ns/op
BenchmarkNonStandardFormat   	221962552	        27.08 ns/op
BenchmarkRFC3339InLocation   	223855030	        26.70 ns/op

// stdlib
BenchmarkGoMultiFormat       	 7471774	       799.0 ns/op
BenchmarkGoRFC3339Nano       	83486947	        69.18 ns/op
BenchmarkGoRFC3339           	100000000	        51.62 ns/op
BenchmarkGoDateTime          	42266613	       140.2 ns/op
BenchmarkGoDateOnly          	70607433	        83.01 ns/op
BenchmarkGoNonStandardFormat 	25277586	       240.5 ns/op
BenchmarkGoRFC3339InLocation 	100000000	        51.40 ns/op
```

# Getting Started

## Installation

```shell
go get github.com/richardliao/parsetime
```

## Usage

```go
package main

import (
	"github.com/richardliao/parsetime"
	"time"
)

func main() {
	// Parse the time to local.
	parsetime.Parse("2006-01-02T15:04:05.999999999+08:00")

	// Parse the time to the specified location.
	loc := time.UTC
	_, locOffset := time.Now().In(loc).Zone()
	parsetime.ParseInLocation("2006-01-02 15:04:05.999", loc, locOffset)
}
```

## Supported Format

ParseTime supports the following formats, as well as their various variants.

```
2006-01-02T15:04:05.999999999+08:00
2006-01-02T15:04:05.999999+08:00
2006-01-02T15:04:05.999+08:00
2006-01-02T15:04:05,999+08:00
2006-01-02T15:04:05+08:00
2006-01-02T15:04:05+0800
2006-01-02T15:04:05+08
2006-01-02T15:04:05Z
2006-01-02T15:04:05z
2006-01-02T15:04:05
2006-01-02 15:04:05
2006-01-02
```
