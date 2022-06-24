package timecode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRate(t *testing.T) {
	t.Parallel()

	rate, err := NewRate(30, false)
	assert.Nil(t, err)
	assert.Equal(t, rate.FPS(), 30.0)
	assert.Equal(t, rate.DropFrame(), false)

	rate, err = NewRate(30, true)
	assert.Nil(t, err)
	assert.Equal(t, rate.FPS(), 30.0)
	assert.Equal(t, rate.DropFrame(), true)

	_, err = NewRate(0, true)
	assert.NotNil(t, err)

	_, err = NewRate(-100, true)
	assert.NotNil(t, err)
}

func TestSMTPERates(t *testing.T) {
	assert.Equal(t, R2997.FPS(), 29.97)
	assert.Equal(t, R2997.DropFrame(), false)

	assert.Equal(t, R2997DF.FPS(), 29.97)
	assert.Equal(t, R2997DF.DropFrame(), true)

	assert.Equal(t, R30.FPS(), 30.0)
	assert.Equal(t, R30.DropFrame(), false)

	assert.Equal(t, R5994.FPS(), 59.94)
	assert.Equal(t, R5994.DropFrame(), false)

	assert.Equal(t, R5994DF.FPS(), 59.94)
	assert.Equal(t, R5994DF.DropFrame(), true)

	assert.Equal(t, R25.FPS(), 25.0)
	assert.Equal(t, R25.DropFrame(), false)

	assert.Equal(t, R50.FPS(), 50.0)
	assert.Equal(t, R50.DropFrame(), false)

	assert.Equal(t, R2398.FPS(), 23.98)
	assert.Equal(t, R2398.DropFrame(), false)

	assert.Equal(t, R60.FPS(), 60.0)
	assert.Equal(t, R60.DropFrame(), false)

	assert.Equal(t, R120.FPS(), 120.0)
	assert.Equal(t, R120.DropFrame(), false)

	assert.Equal(t, R240.FPS(), 240.0)
	assert.Equal(t, R240.DropFrame(), false)
}

func TestParseRate(t *testing.T) {
	rate, err := ParseRate("30000/1001", false)
	assert.Nil(t, err)
	assert.Equal(t, 29.97, rate.FPS())

	rate, err = ParseRate("24000/1001", false)
	assert.Nil(t, err)
	assert.Equal(t, 23.98, rate.FPS())

	rate, err = ParseRate("60000/1001", false)
	assert.Nil(t, err)
	assert.Equal(t, 59.94, rate.FPS())

	rate, err = ParseRate("30/1", false)
	assert.Nil(t, err)
	assert.Equal(t, 30.0, rate.FPS())

	_, err = ParseRate("30/0", false)
	assert.NotNil(t, err)

	_, err = ParseRate("invalid", false)
	assert.NotNil(t, err)

	_, err = ParseRate("", false)
	assert.NotNil(t, err)
}
