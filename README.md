[![Build Status](https://github.com/agorman/go-timecode/workflows/go-timecode-ci/badge.svg)](https://github.com/agorman/go-timecode/actions)
[![go report card](https://goreportcard.com/badge/github.com/agorman/go-timecode "go report card")](https://goreportcard.com/report/github.com/agorman/go-timecode)
[![GoDoc](https://godoc.org/github.com/agorman/go-timecode?status.svg)](https://godoc.org/github.com/agorman/go-timecode)
[![codecov](https://codecov.io/gh/agorman/go-timecode/branch/master/graph/badge.svg)](https://codecov.io/gh/agorman/go-timecode)

A package for dealing with SMPTE timecode.

It is primarily based on https://www.npmjs.com/package/timecode

```
frameRate := 29.97
dropFrame := true

tc := timecode.NewTimecode(frameRate, dropFrame)
tc.AddSeconds(10.5)
tc.AddFrames(20)
tc.AddString("05:01:20;18")
tc.SubSeconds(1)

tc2 := timecode.NewTimecode(frameRate, dropFrame)
tc.AddString("05:01:20;18")

tc.Add(tc2)

fmt.Println(tc) // 05:01:30;23
```
