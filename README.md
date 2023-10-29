# Overview

ParseTime is a Go package dedicated to quickly parse RFC3339 or similar time strings.

The time.Parse in Go stdlib is fast enough in most cases, but not when processing massive time series data.
The time string generated in many languages is not the standard RFC3339, but a variant of it.
When it is not a standard RFC3339 string (e.g. 2006-01-02T15:04:05.999), 
or multiple uncertain time formats need to be tried, 
the resources consumed by the parsing time are very considerable.


```
// parsetime
BenchmarkRFC3339Nano    58373623                55.30 ns/op

// stdlib
BenchmarkGoRFC3339Nano  46511388                76.82 ns/op
```
