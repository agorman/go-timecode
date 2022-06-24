package timecode

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

var (
	timecodeRegExp = regexp.MustCompile(`^(\d+)[:;.,](\d\d)[:;.,](\d\d)[:;.,](\d+)$`)
)

// Timecode is used to simplify using string based timecodes by providing conversions, frame based math,
// and support for SMTPE drop frame encoding. This timecode library supports hours of any length and does
// not loop back to 00:00:00:00 after 59:59:59:{fps-1}.
type Timecode struct {
	rate   Rate
	frames uint64
}

// Parse takes rate and a timecode as a string in the form hh:mm:ss:ff. Where hh represents hours, mm represents minutes,
// ss represents seconds, and ff represents frames. Minutes and seconds must between 0 and 59. Hours and frames must be
// greather than or equal to 0. Hours, minutes, seconds, or frames less than 10 must be left padded with a 0. This means
// that negative timecodes are not supported. The separator isn't required to be : and will match any of [:;,.] in any position.
// Parse is written to be as forgiving as possible.
func Parse(rate Rate, s string) (Timecode, error) {
	tc := Timecode{
		rate: rate,
	}

	matches := timecodeRegExp.FindStringSubmatch(s)
	if len(matches) != 5 {
		return tc, fmt.Errorf("unable to parse timecode: %s", s)
	}

	hours, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return tc, fmt.Errorf("unable to parse timecode hours: %s: %w", s, err)
	}

	minutes, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return tc, fmt.Errorf("unable to parse timecode minutes: %s: %w", s, err)
	}
	if minutes >= 60 {
		return tc, fmt.Errorf("minutes must be between 0 and 59 got: %d", uint64(minutes))
	}

	seconds, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return tc, fmt.Errorf("unable to parse timecode minutes: %s: %w", s, err)
	}
	if seconds >= 60 {
		return tc, fmt.Errorf("minutes must be between 0 and 59 got: %d", uint64(seconds))
	}

	frames, err := strconv.ParseFloat(matches[4], 64)
	if err != nil {
		return tc, fmt.Errorf("unable to parse timecode minutes: %s: %w", s, err)
	}
	if frames >= tc.rate.timeBase {
		return tc, fmt.Errorf("frames must be between 0 and %f got: %d", tc.rate.fps, uint64(frames))
	}

	hourFrames := rate.timeBase * 3600
	minutesFrames := rate.timeBase * 60
	totalFrames := hourFrames*hours + minutesFrames*minutes + rate.timeBase*seconds + frames

	if rate.dropFrame {
		dropFrames := math.Round(rate.fps * 0.066666)
		totalMinutes := (60 * hours) + minutes
		// remove skipped frames from the frame count
		frames := totalFrames - (dropFrames * (totalMinutes - totalMinutes/10))
		tc.frames = uint64(frames)
	} else {
		tc.frames = uint64(totalFrames)
	}

	return tc, nil
}

// FromFrames returns a Timecode based on the passed rate and frames.
func FromFrames(rate Rate, frames uint64) Timecode {
	return Timecode{
		rate:   rate,
		frames: frames,
	}
}

// FromSeconds returns a Timecode based on the passed rate and seconds.
func FromSeconds(rate Rate, seconds float64) (Timecode, error) {
	tc := Timecode{
		rate: rate,
	}

	if seconds < 0 {
		return tc, fmt.Errorf("timecode can not have a negative value: %f", seconds)
	}

	totalFrames := seconds * rate.timeBase
	tc.frames = uint64(totalFrames)

	return tc, nil
}

// Rate returns the Rate used when creating the Timecode.
func (tc Timecode) Rate() Rate {
	return tc.rate
}

// Hour returns the hour portion of the timecode string as a uint64.
// For example a timecode of 02:12:49:15 would return 2.
func (tc Timecode) Hour() uint64 {
	var hour uint64
	if tc.rate.dropFrame {
		hour, _, _, _ = tc.dropFrameToParts()
	} else {
		hour, _, _, _ = tc.toParts()
	}

	return hour
}

// Minute returns the minute portion of the timecode string as a uint64.
// For example a timecode of 02:12:49:15 would return 12.
func (tc Timecode) Minute() uint64 {
	var minute uint64
	if tc.rate.dropFrame {
		_, minute, _, _ = tc.dropFrameToParts()
	} else {
		_, minute, _, _ = tc.toParts()
	}

	return minute
}

// Second returns the second portion of the timecode string as a uint64.
// For example a timecode of 02:12:49:15 would return 49.
func (tc Timecode) Second() uint64 {
	var second uint64
	if tc.rate.dropFrame {
		_, _, second, _ = tc.dropFrameToParts()
	} else {
		_, _, second, _ = tc.toParts()
	}

	return second
}

// Frame returns the frame portion of the timecode string as a uint64.
// For example a timecode of 02:12:49:15 would return 15.
func (tc Timecode) Frame() uint64 {
	var frame uint64
	if tc.rate.dropFrame {
		_, _, _, frame = tc.dropFrameToParts()
	} else {
		_, _, _, frame = tc.toParts()
	}

	return frame
}

// String returns the entire timecode formatted as a string based on the frame rate and drop frame
// encoding.
func (tc Timecode) String() string {
	var hour, minute, second, frame uint64
	if tc.rate.dropFrame {
		hour, minute, second, frame = tc.dropFrameToParts()
	} else {
		hour, minute, second, frame = tc.toParts()
	}

	sep := ":"
	if tc.rate.dropFrame {
		sep = ";"
	}

	return fmt.Sprintf("%02d:%02d:%02d%s%02d", hour, minute, second, sep, frame)
}

// Frames returns the frames as an int64 based on the frame rate and drop frame
// encoding.
func (tc Timecode) Frames() uint64 {
	return tc.frames
}

// Seconds returns the total seconds as an float64 based on the frame rate and drop frame
// encoding.
func (tc Timecode) Seconds() float64 {
	hour, minute, second, frame := tc.toParts()
	return float64(hour*3600+minute*60+second) + float64(frame)/tc.rate.timeBase
}

// Add adds the frames to the Timecode and returns a new Timecode as the result.
func (tc Timecode) Add(frames uint64) Timecode {
	return FromFrames(tc.rate, tc.Frames()+frames)
}

// Sub subtracts the frames from the Timecode and returns a new Timecode as the result.
// If the result would be a negative timecode then an error is returned.
func (tc Timecode) Sub(frames uint64) (Timecode, error) {
	result := int64(tc.Frames() - frames)
	if result < 0 {
		return Timecode{}, fmt.Errorf("resulting timecode would have a negative value: %d", result)
	}

	return FromFrames(tc.rate, uint64(result)), nil
}

func (tc Timecode) dropFrameToParts() (uint64, uint64, uint64, uint64) {
	dropFrames := uint64(math.Round(tc.rate.fps * 0.066666))

	// framesPerHour := uint64(math.Round(tc.rate.fps * 3600))
	// framesPer24Hours := framesPerHour * 24
	framesPer10Minutes := uint64(math.Round(tc.rate.fps * 60 * 10))
	framesPerMinute := (uint64(math.Round(tc.rate.fps)) * 60) - dropFrames

	framenumber := tc.frames
	d := framenumber / framesPer10Minutes
	m := framenumber % framesPer10Minutes

	if m > dropFrames {
		framenumber = framenumber + (dropFrames * 9 * d) + dropFrames*((m-dropFrames)/framesPerMinute)
	} else {
		framenumber = framenumber + dropFrames*9*d
	}

	// frRound = math.Round(framerate);
	frame := framenumber % uint64(tc.rate.timeBase)
	second := (framenumber / uint64(tc.rate.timeBase)) % 60
	minute := ((framenumber / uint64(tc.rate.timeBase)) / 60) % 60
	hour := (((framenumber / uint64(tc.rate.timeBase)) / 60) / 60)

	return hour, minute, second, frame
}

func (tc Timecode) toParts() (uint64, uint64, uint64, uint64) {
	framesPerHour := uint64(tc.rate.timeBase) * 3600
	framesPerMinute := uint64(tc.rate.timeBase) * 60

	remaining := tc.frames

	hour := remaining / framesPerHour
	remaining = remaining - (hour * framesPerHour)

	minute := remaining / framesPerMinute
	remaining = remaining - (minute * framesPerMinute)

	second := remaining / uint64(tc.rate.timeBase)
	remaining = remaining - (second * uint64(tc.rate.timeBase))

	return uint64(hour), uint64(minute), uint64(second), uint64(remaining)
}
