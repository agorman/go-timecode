[![Build Status](https://github.com/agorman/go-timecode/workflows/go-timecode-ci/badge.svg)](https://github.com/agorman/go-timecode/actions)
[![go report card](https://goreportcard.com/badge/github.com/agorman/go-timecode "go report card")](https://goreportcard.com/report/github.com/agorman/go-timecode)
[![GoDoc](https://godoc.org/github.com/agorman/go-timecode/v2?status.svg)](https://godoc.org/github.com/agorman/go-timecode/v2)
[![codecov](https://codecov.io/gh/agorman/go-timecode/branch/master/graph/badge.svg)](https://codecov.io/gh/agorman/go-timecode)

# go-timecode

go-timecode simplifies the use of string based timecodes by providing conversions, frame based math, and support for SMTPE drop frame encoding. go-timecode  offers a variety of industry standard formats out of the box but is designed to make it easy to work with any combination of formats you
need.

## Installation

```
go get github.com/agorman/go-timecode/v2

```

## Documentation

https://godoc.org/github.com/agorman/go-timecode/v2

## Basic usage

~~~
tc, err := timecode.Parse(timecode.R30, "01:30:12:15")
if err != nil {
    panic(err)
}

tc.String()   # "01:30:12:15"
tc.Frames()   # 162375
tc.Seconds()  # 5412.5
~~~

~~~
tc := timecode.FromFrames(timecode.R2997DF, 162213)

tc.String()   # "01:30:12:15"
tc.Frames()   # 162213
tc.Seconds()  # 5407.1
~~~

~~~
tc, err := timecode.FromSeconds(timecode.R2398, 5412.625)
if err != nil {
    panic(err)
}

tc.String()   # "01:30:12:15"
tc.Frames()   # 129903
tc.Seconds()  # 5412.625
~~~

~~~
rate, err := timecode.NewRate(30, false)
if err != nil {
    panic(err)
}
rate.FPS()         # 30.0
rate.DropFrame()   # false
~~~

~~~
rate, err := timecode.ParseRate("30000/1001", true)
if err != nil {
    panic(err)
}
rate.FPS()         # 29.97
rate.DropFrame()   # true
~~~
