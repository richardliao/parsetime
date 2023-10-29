# Overview

ParseTime is a Go package dedicated to quickly parse RFC3339 or similar time strings.

The time.Parse in Go stdlib is fast enough in most cases, but not when processing massive time series data.
The time string generated in many languages is not the standard RFC3339, but a variant of it.
When it is not a standard RFC3339 string (e.g. 2006-01-02T15:04:05.999), 
or multiple uncertain time formats need to be tried, 
the resources consumed by the parsing time are very considerable.

```
// parsetime
BenchmarkRFC3339NanoBytes       123328021               28.68 ns/op
BenchmarkMultiFormat            100000000               33.57 ns/op
BenchmarkRFC3339Nano            100000000               33.40 ns/op
BenchmarkRFC3339                123554547               29.54 ns/op
BenchmarkDateTime               168250800               21.22 ns/op
BenchmarkDateOnly               234963064               15.13 ns/op
BenchmarkNonStandardFormat      100000000               33.35 ns/op
BenchmarkRFC3339InLocation      100000000               33.45 ns/op

// stdlib
BenchmarkGoMultiFormat           4040998               897.9 ns/op
BenchmarkGoRFC3339Nano          44425657                76.66 ns/op
BenchmarkGoRFC3339              62792282                56.75 ns/op
BenchmarkGoDateTime             21754081               165.3 ns/op
BenchmarkGoDateOnly             35960893               101.2 ns/op
BenchmarkGoNonStandardFormat    13469737               271.2 ns/op
BenchmarkGoRFC3339InLocation    63190134                58.61 ns/op
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
