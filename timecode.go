package timecode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Timecode is an object for manipulating SMPTE timecodes
type Timecode struct {
	frameRate    float64
	intFrameRate int
	dropFrame    bool
	frameCount   int
	tcRegexp     *regexp.Regexp
}

// NewTimecode returns a new Timecode.
// It takes the frame rate and a bool which flags it as using drop frame encoding
func NewTimecode(rate float64, drop bool) (*Timecode, error) {
	tcRegexp := regexp.MustCompile(`^(\d\d)[:;](\d\d)[:;](\d\d)[:;](\d+)$`)

	intFrameRate := int(rate + 0.5)
	if intFrameRate <= 0 {
		return nil, fmt.Errorf("Unsupported frame rate")
	}

	return &Timecode{
		frameRate:    rate,
		intFrameRate: intFrameRate,
		dropFrame:    drop,
		frameCount:   0,
		tcRegexp:     tcRegexp,
	}, nil
}

func (tc *Timecode) SetFrameRate(rate float64) error {
	intFrameRate := int(rate + 0.5)
	if intFrameRate <= 0 {
		return fmt.Errorf("Unsupported frame rate")
	}

	tc.frameRate = rate
	tc.intFrameRate = intFrameRate

	return nil
}

func (tc *Timecode) SetDropFrame(drop bool) {
	tc.dropFrame = drop
}

// Resets sets the timecode to 0
func (tc *Timecode) Reset() {
	tc.frameCount = 0

	return
}

// AddString adds SMPTE strings in the following formats: timecode strings non-drop 'hh:mm:ss:ff', drop 'hh:mm:ss;ff', or milliseconds 'hh:mm:ss:mmm'
func (tc *Timecode) AddString(t string) error {
	frames, err := tc.timecodeToFrames(t)
	if err != nil {
		return err
	}

	tc.frameCount += frames

	return nil
}

// AddSeconds add seconds to the timecode
func (tc *Timecode) AddSeconds(seconds float64) *Timecode {
	tc.frameCount += int(float64(tc.intFrameRate)*seconds + 0.5)

	return tc
}

// AddFrames add frames to the timecode
func (tc *Timecode) AddFrames(frames int) *Timecode {
	tc.frameCount += frames

	return tc
}

// Add adds another Timecode object to the timecode
func (tc *Timecode) Add(t *Timecode) *Timecode {
	tc.frameCount += t.frameCount

	return tc
}

// SubString subtract SMPTE strings timecode strings non-drop 'hh:mm:ss:ff', drop 'hh:mm:ss;ff', or milliseconds 'hh:mm:ss:mmm'
func (tc *Timecode) SubString(t string) error {
	frames, err := tc.timecodeToFrames(t)
	if err != nil {
		return err
	}

	tc.frameCount -= frames

	return nil
}

// SubSeconds subtract seconds from the timecode
func (tc *Timecode) SubSeconds(seconds float64) *Timecode {
	tc.frameCount -= int(float64(tc.intFrameRate)*seconds + 0.5)

	return tc
}

// SubFrames subtract frames from the timecode
func (tc *Timecode) SubFrames(frames int) *Timecode {
	tc.frameCount -= frames

	return tc
}

// Subtract subtracts another Timecode object from the timecode
func (tc *Timecode) Sub(t *Timecode) *Timecode {
	tc.frameCount -= t.frameCount

	return tc
}

// String returns a SMTPE timecode string
func (tc *Timecode) String() string {
	return tc.framesToTimecode(tc.frameCount)
}

func (tc *Timecode) framesToTimecode(frameCount int) string {
	if tc.dropFrame {
		dropFrames := int(tc.frameRate*0.066666 + 0.5)
		framesPerHour := int(tc.frameRate*60*60 + 0.5)
		framesPer24Hours := framesPerHour * 24
		framesPer10Minutes := int(tc.frameRate*60*10 + 0.5)
		framesPerMinute := int(tc.frameRate*60 + 0.5)

		// roll over clock if greater than 24 hours
		frameCount = frameCount % framesPer24Hours

		// if time is negative, count back from 24 hours
		if frameCount < 0 {
			frameCount = framesPer24Hours + frameCount
		}

		d := int(frameCount / framesPer10Minutes)
		m := frameCount % framesPer10Minutes

		if m > dropFrames {
			frameCount = frameCount + (dropFrames * 9 * d) + dropFrames*int((m-dropFrames)/framesPerMinute)
		} else {
			frameCount = frameCount + dropFrames*9*d
		}

		return fmt.Sprintf("%02d:%02d:%02d;%02d",
			int(int(int(frameCount/tc.intFrameRate)/60)/60),
			int(int(frameCount/tc.intFrameRate)/60)%60,
			int(frameCount/tc.intFrameRate)%60,
			frameCount%tc.intFrameRate,
		)
	}

	hours := tc.frameCount / (3600 * tc.intFrameRate)
	if hours > 23 {
		hours = hours % 24
		frameCount = frameCount - (23 * 3600 * tc.intFrameRate)
	}

	minutes := (frameCount % (3600 * tc.intFrameRate)) / (60 * tc.intFrameRate)
	seconds := ((frameCount % (3600 * tc.intFrameRate)) % (60 * tc.intFrameRate)) / tc.intFrameRate
	frames := ((frameCount % (3600 * tc.intFrameRate)) % (60 * tc.intFrameRate)) % tc.intFrameRate

	return fmt.Sprintf("%02d:%02d:%02d:%02d", hours, minutes, seconds, frames)
}

func (tc *Timecode) timecodeToFrames(t string) (int, error) {
	// parses timecode strings non-drop 'hh:mm:ss:ff', drop 'hh:mm:ss;ff', or milliseconds 'hh:mm:ss:fff'
	var hours, minutes, seconds, frames, ms, totalSeconds float64

	if !tc.tcRegexp.MatchString(t) {
		return 0, fmt.Errorf("Timecode string parsing error. %s", t)
	}

	if len(t) == 11 {
		frames, _ = strconv.ParseFloat(t[9:11], 32)
		ms = frames / float64(tc.intFrameRate)
	} else if len(t) == 12 {
		ms, _ = strconv.ParseFloat(t[9:], 32)
		ms /= 1000
	} else {
		return 0, fmt.Errorf("Timecode string parsing error. %s", t)
	}

	// these are correct because of matching regexp
	hours, _ = strconv.ParseFloat(t[0:2], 32)
	minutes, _ = strconv.ParseFloat(t[3:5], 32)
	seconds, _ = strconv.ParseFloat(t[6:8], 32)

	if strings.Contains(t, `;`) {
		dropFrames := int(tc.frameRate*0.066666 + 0.5)
		hourFrames := tc.intFrameRate * 60 * 60
		minuteFrames := tc.intFrameRate * 60
		totalMinutes := (hours * 60) + minutes
		return ((hourFrames * int(hours)) + (minuteFrames * int(minutes)) + (tc.intFrameRate * int(seconds)) + int(frames)) - (dropFrames * (int(totalMinutes) - int(int(totalMinutes)/10))), nil
	}

	totalSeconds = hours*3600 + minutes*60 + seconds + ms
	return int(float64(tc.intFrameRate)*totalSeconds + 0.5), nil
}
