package timecode

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

var (
	rateRegExp = regexp.MustCompile(`^(\d+)/(\d+)$`)
)

var (
	// R2297 is SMPTE (Society of Motion Picture and Television Engineers) 29.97 fps NDF
	R2997 = Rate{fps: 29.97, timeBase: 30, dropFrame: false}

	// R2997DF is SMPTE (Society of Motion Picture and Television Engineers) 29.97 fps DF
	R2997DF = Rate{fps: 29.97, timeBase: 30, dropFrame: true}

	// R30 is SMPTE (Society of Motion Picture and Television Engineers) 30 fps NDF
	R30 = Rate{fps: 30, timeBase: 30, dropFrame: false}

	// R5994 is SMPTE (Society of Motion Picture and Television Engineers) 59.94 fps NDF
	R5994 = Rate{fps: 59.94, timeBase: 60, dropFrame: false}

	// R5994DF is SMPTE (Society of Motion Picture and Television Engineers) 59.94 fps DF
	R5994DF = Rate{fps: 59.94, timeBase: 60, dropFrame: true}

	// R60 is SMTPE (Society of Motion Picture and Television Engineers) 60 fps NDF
	R60 = Rate{fps: 60, timeBase: 60, dropFrame: false}

	// R25 is EBU (European Broadcasting Union 25 fps NDF
	R25 = Rate{fps: 25, timeBase: 25, dropFrame: false}

	// R50 is EBU (European Broadcasting Union 50 fps NDF
	R50 = Rate{fps: 50, timeBase: 50, dropFrame: false}

	// R2398 is Film 23.99 fps NDF
	R2398 = Rate{fps: 23.98, timeBase: 24, dropFrame: false}

	// R24 is Film 23.99 fps NDF
	R24 = Rate{fps: 24.0, timeBase: 24, dropFrame: false}

	// R120 is Film 120 fps NDF
	R120 = Rate{fps: 120, timeBase: 120, dropFrame: false}

	// R240 is SMTPE 240 fps NDF
	R240 = Rate{fps: 240, timeBase: 240, dropFrame: false}
)

// Rate describes a frame rate and drop frame encoding for a Timecode.
type Rate struct {
	fps       float64
	timeBase  float64
	dropFrame bool
}

// NewRate returns a Rate baed on the given fps (frame rate) and dropFrame. The
// fps must be an integer greater than or equal to 1.  The resulting FPS is rounded
// to two decimal places. A drop frame timecode is a SMPTE standard that works by
// skipping two frames per minute except for every 10th minute.
func NewRate(fps float64, dropFrame bool) (Rate, error) {
	rate := Rate{}

	fps = math.Round(fps*100) / 100

	if fps < 1 {
		return rate, fmt.Errorf("rate must be at least 1 fps but got: %f", fps)
	}

	timeBase := math.Round(fps)

	return Rate{
		fps:       fps,
		timeBase:  timeBase,
		dropFrame: dropFrame,
	}, nil
}

// ParseRate takes a frame rate string in a fractional form and returns a Rate object.
// The string is of the form num/den where num is an integer that's greater than 0 and
// den is an integer greater than 1. The resulting FPS is rounded to two decimal places.
func ParseRate(s string, dropFrame bool) (Rate, error) {
	rate := Rate{
		dropFrame: dropFrame,
	}

	matches := rateRegExp.FindStringSubmatch(s)
	if len(matches) != 3 {
		return rate, fmt.Errorf("unable to parse rate: %s", s)
	}

	num, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return rate, fmt.Errorf("unable to parse rate numerator: %s: %w", s, err)
	}

	den, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return rate, fmt.Errorf("unable to parse rate denominator: %s: %w", s, err)
	}
	if den == 0 {
		return rate, fmt.Errorf("rate cannot have a denominator of 0: %s", s)
	}

	// round to two decimal places
	fps := num / den
	if fps < 1 {
		return rate, fmt.Errorf("rate must be at least 1 fps but got: %f", fps)
	}

	rate.fps = math.Round(fps*100) / 100
	return rate, nil

}

// FPS returns the fps (frame rate).
func (r Rate) FPS() float64 {
	return r.fps
}

// DropFrame returns true if this Rate is using SMTPE drop frame encoding.
func (r Rate) DropFrame() bool {
	return r.dropFrame
}
