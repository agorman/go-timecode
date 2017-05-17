A package for dealing with SMPTE timecode.

It is primarily based on https://www.npmjs.com/package/timecode

```
frameRate := 29.97
dropFrame := true

tc := timecode.NewTimecode(frameRate, dropFrame)
tc.AddSeconds(10.5)
tc.AddFrames(20)
tc.SubSeconds(1)

fmt.Println(tc) // 00:00:10;05
```
