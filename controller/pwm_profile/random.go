package pwm_profile

import (
	"encoding/json"
	"math/rand"
	"time"
)

const (
	seed         = 154
	peakInterval = 160
)

type random struct {
	temporal
	previous     float64
	peaks        []float64
	peakInterval float64
}

func (r *random) Name() string {
	return _randomProfileName
}

func Random(conf json.RawMessage, min, max float64) (*random, error) {
	t, err := Temporal(conf, min, max)
	if err != nil {
		return nil, err
	}
	return NewRandom(t), nil
}

func NewRandom(t temporal) *random {
	rand.Seed(seed)
	numPeaks := int(t.TotalSeconds() / peakInterval)
	if numPeaks == 0 {
		numPeaks = 1
	}
	peaks := make([]float64, numPeaks)
	for i, _ := range peaks {
		peaks[i] = rand.Float64()*t.ValueRange() + t.min
	}
	return &random{
		temporal:     t,
		previous:     peaks[0],
		peakInterval: peakInterval,
		peaks:        peaks,
	}
}

func (s *random) Get(t time.Time) float64 {
	if s.IsOutside(t) {
		return 0
	}
	i := int(s.PastSeconds(t) / s.peakInterval)
	if i >= len(s.peaks) {
		i = len(s.peaks) - 1
	}
	prevPeak := s.peaks[i]
	nextPeak := s.peaks[0]
	if (i + 1) < len(s.peaks) {
		nextPeak = s.peaks[i+1]
	}
	f := (nextPeak - prevPeak) / s.peakInterval
	inc := rand.NormFloat64() + f
	s.previous += inc
	if s.previous > s.max {
		s.previous = s.max
	}
	if s.previous < s.min {
		s.previous = s.min
	}
	s.previous = float64(int(s.previous*100)) / 100
	return s.previous
}
