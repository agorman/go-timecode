package timecode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrameRateTimecode(t *testing.T) {
	t.Parallel()

	_, err := NewTimecode(-1, false)
	assert.NotNil(t, err)

	tc, err := NewTimecode(30, false)
	assert.Nil(t, err)

	err = tc.SetFrameRate(-5)
	assert.NotNil(t, err)

	tc, err = NewTimecode(30, false)
	assert.Nil(t, err)

	err = tc.SetFrameRate(0)
	assert.NotNil(t, err)
}

func TestResetTimecode(t *testing.T) {
	t.Parallel()

	tc, err := NewTimecode(30, false)
	assert.Nil(t, err)

	tc.AddSeconds(90)
	assert.Equal(t, 90.0, tc.ToSeconds())

	tc.Reset()
	assert.Equal(t, 0.0, tc.ToSeconds())
	assert.Equal(t, "00:00:00:00", tc.String())
}

func TestStringTimecode(t *testing.T) {
	t.Parallel()

	tc, err := NewTimecode(30, false)
	assert.Nil(t, err)

	err = tc.AddString("00:asdf:123")
	assert.NotNil(t, err)

	err = tc.SubString("1:2:3:4")
	assert.NotNil(t, err)
}

func TestAddTimecode(t *testing.T) {
	t.Parallel()

	frameRate := 30.0
	dropFrame := false

	tc, err := NewTimecode(frameRate, dropFrame)
	assert.Nil(t, err)
	assert.Equal(t, "00:00:00:00", tc.String())
	assert.Equal(t, 0.0, tc.ToSeconds())

	tc.AddSeconds(90)
	assert.Equal(t, "00:01:30:00", tc.String())
	assert.Equal(t, 90.0, tc.ToSeconds())

	tc.AddSeconds(3760)
	assert.Equal(t, "01:04:10:00", tc.String())
	assert.Equal(t, 3850.0, tc.ToSeconds())

	tc.AddFrames(15)
	assert.Equal(t, "01:04:10:15", tc.String())
	assert.Equal(t, 3850.5, tc.ToSeconds())

	tc.AddFrames(100)
	assert.Equal(t, "01:04:13:25", tc.String())
	assert.Equal(t, 3853.8333333333335, tc.ToSeconds())

	err = tc.AddString("01:02:03:04")
	assert.Nil(t, err)
	assert.Equal(t, "02:06:16:29", tc.String())
	assert.Equal(t, 7576.966666666666, tc.ToSeconds())

	err = tc.AddString("00:55:45:02")
	assert.Nil(t, err)
	assert.Equal(t, "03:02:02:01", tc.String())
	assert.Equal(t, 10922.033333333333, tc.ToSeconds())

	frameRate2 := 25.0
	dropFrame2 := false
	tc2, err := NewTimecode(frameRate2, dropFrame2)
	assert.Nil(t, err)
	tc2.AddSeconds(100)
	assert.Equal(t, "00:01:40:00", tc2.String())
	assert.Equal(t, 100.0, tc2.ToSeconds())

	tc.Add(tc2)
	assert.Equal(t, 11005.366666666667, tc.ToSeconds())
}

func TestSubTimecode(t *testing.T) {
	t.Parallel()

	frameRate := 30.0
	dropFrame := false

	tc, err := NewTimecode(frameRate, dropFrame)
	assert.Nil(t, err)

	err = tc.AddString("01:01:01:00")
	assert.Nil(t, err)
	assert.Equal(t, "01:01:01:00", tc.String())
	assert.Equal(t, 3661.0, tc.ToSeconds())

	tc.SubSeconds(90)
	assert.Equal(t, "00:59:31:00", tc.String())
	assert.Equal(t, 3571.0, tc.ToSeconds())

	tc.SubSeconds(760)
	assert.Equal(t, "00:46:51:00", tc.String())
	assert.Equal(t, 2811.0, tc.ToSeconds())

	tc.SubFrames(15)
	assert.Equal(t, "00:46:50:15", tc.String())
	assert.Equal(t, 2810.5, tc.ToSeconds())

	tc.SubFrames(100)
	assert.Equal(t, "00:46:47:05", tc.String())
	assert.Equal(t, 2807.1666666666665, tc.ToSeconds())

	err = tc.SubString("00:30:03:04")
	assert.Nil(t, err)
	assert.Equal(t, "00:16:44:01", tc.String())
	assert.Equal(t, 1004.0333333333333, tc.ToSeconds())

	frameRate2 := 25.0
	dropFrame2 := false
	tc2, err := NewTimecode(frameRate2, dropFrame2)
	assert.Nil(t, err)

	tc.Sub(tc2)
	assert.Equal(t, 1004.0333333333333, tc.ToSeconds())
}

func TestDropFrameTimecode(t *testing.T) {
	t.Parallel()

	frameRate := 29.97
	dropFrame := true

	tc, err := NewTimecode(frameRate, dropFrame)
	assert.Nil(t, err)
	assert.Equal(t, "00:00:00;00", tc.String())
	assert.Equal(t, 0.0, tc.ToSeconds())

	tc.AddSeconds(90)
	assert.Equal(t, "00:01:30;02", tc.String())
	assert.Equal(t, 90.09009009009009, tc.ToSeconds())

	tc.SubSeconds(89)
	assert.Equal(t, "00:00:01;00", tc.String())
	assert.Equal(t, 1.001001001001001, tc.ToSeconds())

	tc.SetDropFrame(false)
	assert.Equal(t, "00:00:01:00", tc.String())
	assert.Equal(t, 1.001001001001001, tc.ToSeconds())

	tc.SetDropFrame(true)
	assert.Equal(t, "00:00:01;00", tc.String())
	assert.Equal(t, 1.001001001001001, tc.ToSeconds())
}
