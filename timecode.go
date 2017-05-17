package timecode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/*
 * Timecode object for manipulating SMPTE timecodes
 */
type Timecode struct {
	frameRate    float64
	intFrameRate int
	dropFrame    bool
	frameCount   int
	tcRegexp     *regexp.Regexp
}

/*
 * Create a new timecode
 *
 * Takes the frame rate and a bool which flags it as using drop frame encoding
 */
func NewTimecode(rate float64, drop bool) *Timecode {
	tcRegexp := regexp.MustCompile(`^(\d\d)[:;](\d\d)[:;](\d\d)[:;](\d+)$`)

	return &Timecode{
		frameRate:    rate,
		intFrameRate: int(rate + 0.5),
		dropFrame:    drop,
		frameCount:   0,
		tcRegexp:     tcRegexp,
	}
}

/*
 * Resets the timecode to 0
 */
func (tc *Timecode) Reset() {
	tc.frameCount = 0

	return
}

/*
 * Add a SMPTE timecode to the timecode
 *
 * Takes timecode strings non-drop 'hh:mm:ss:ff', drop 'hh:mm:ss;ff', or milliseconds 'hh:mm:ss:mmm'
 */
func (tc *Timecode) AddString(t string) error {
	frames, err := tc.timecodeToFrames(t)
	if err != nil {
		return err
	}

	tc.frameCount += frames

	return nil
}

/*
 * Add seconds to the timecode
 */
func (tc *Timecode) AddSeconds(seconds float64) {
	tc.frameCount += int(float64(tc.intFrameRate)*seconds + 0.5)
}

/*
 * Add frames to the timecode
 */
func (tc *Timecode) AddFrames(frames int) {
	tc.frameCount += frames
}

/*
 * Add another Timecode object to the timecode
 */
func (tc *Timecode) Add(t *Timecode) {
	tc.frameCount += t.frameCount
}

/*
 * Subtract a SMPTE timecode from the timecode
 *
 * Takes timecode strings non-drop 'hh:mm:ss:ff', drop 'hh:mm:ss;ff', or milliseconds 'hh:mm:ss:mmm'
 */
func (tc *Timecode) SubString(t string) error {
	frames, err := tc.timecodeToFrames(t)
	if err != nil {
		return err
	}

	tc.frameCount -= frames

	return nil
}

/*
 * Subtract seconds from the timecode
 */
func (tc *Timecode) SubSeconds(seconds float64) {
	tc.frameCount -= int(float64(tc.intFrameRate)*seconds + 0.5)
}

/*
 * Subtract frames from the timecode
 */
func (tc *Timecode) SubFrames(frames int) {
	tc.frameCount -= frames
}

/*
 * Subtract another Timecode object from the timecode
 */
func (tc *Timecode) Sub(t *Timecode) {
	tc.frameCount -= t.frameCount
}

/*
 * Return a SMTPE timecode string
 */
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
