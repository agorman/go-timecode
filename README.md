A package for dealing with SMPTE timecode.

It is primarily based on https://www.npmjs.com/package/timecode

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](http://godoc.org/github.com/agorman/go-timecode)

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
