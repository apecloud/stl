package stl

import (
	"github.com/apecloud/stl/loess"
)

// Config is a configuration structure
type Config struct {
	// Width represents the width of the LOESS smoother in data points
	Width int

	// Jump is the number of points to skip between smoothing,
	Jump int

	// Which weight updating function should be used?
	Fn loess.WeightUpdate
}

// Opt is a function that helps build the conf
type Opt func(*state)

// WithSeasonalConfig configures the seasonal component of the decomposition process.
func WithSeasonalConfig(conf Config) Opt {
	return func(s *state) {
		s.sConf = conf
	}
}

// WithTrendConfig configures the trend component of the decomposition process.
func WithTrendConfig(conf Config) Opt {
	return func(s *state) {
		s.tConf = conf
	}
}

// WithLowpassConfig configures the operation that performs the lowpass filter in the decomposition process
func WithLowpassConfig(conf Config) Opt {
	return func(s *state) {
		s.lConf = conf
	}
}

// WithRobustIter indicates how many iterations of "robust" (i.e. outlier removal) to do.
// The default is 0.
func WithRobustIter(n int) Opt {
	return func(s *state) {
		s.robustIter = n
	}
}

// WithIter indicates how many iterations to run.
// The default is 2.
func WithIter(n int) Opt {
	return func(s *state) {
		s.innerIter = n
	}
}

// WithQuadratic indicates how many iterations to run.
// The default is 2.
func WithQuadratic() Opt {
	return func(s *state) {
		s.tConf.Fn = loess.Quadratic
	}
}

// DefaultSeasonal returns the default configuration for the operation that works on the seasonal component.
func DefaultSeasonal(width int) Config {
	if width <= 0 {
		panic("Cannot use negative window width for seasonal smoothing.")
	}
	jmp := int(0.1 * float64(width))
	if jmp <= 0 {
		jmp = 1
	}
	return Config{
		Width: width,
		Jump:  jmp,
		Fn:    loess.Linear,
	}
}

// DefaultTrend returns the default configuration for the operation that works on the trend component.
func DefaultTrend(periodicity, width int) Config {
	if periodicity <= 0 {
		panic("Cannot use negative periodicity for trend smoothing.")
	}
	if width <= 0 {
		panic("Cannot use negative window width for trend smoothing.")
	}
	jmp := int(0.1 * float64(width))
	if jmp <= 0 {
		jmp = 1
	}

	return Config{
		Width: trendWidth(periodicity, width),
		Jump:  jmp,
		Fn:    loess.Linear,
	}
}

// DefaultLowPass returns the default configuration for the operation that works on the lowpass component.
func DefaultLowPass(periodicity int) Config {
	if periodicity <= 0 {
		panic("Cannot use negative periodicity for low pass smoothing.")
	}
	jmp := int(0.1 * float64(periodicity))
	if jmp <= 0 {
		jmp = 1
	}
	return Config{
		Width: periodicity,
		Jump:  jmp,
		Fn:    loess.Linear,
	}
}

// from the original paper's numerical analysis section
func trendWidth(periodicity int, seasonalWidth int) int {
	p := float64(periodicity)
	w := float64(seasonalWidth)

	de := int(1.5*p/(1-1.5/w) + 0.5)

	if de <= 0 {
		panic("Default trend width calculation overflowed")
	}
	return de
}
