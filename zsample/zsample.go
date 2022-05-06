package zsample

import (
	"time"
)

func NewFloatSampler(size int, frequency time.Duration) *FloatSampler {
	return &FloatSampler{
		size:      size,
		frequency: frequency,
		samples:   make([]float64, size),
		times:     make([]time.Time, size),
	}
}

type FloatSampler struct {
	size      int
	frequency time.Duration
	index     int
	samples   []float64
	times     []time.Time
}

func (s *FloatSampler) Add(val float64) (added bool) {
	n := time.Now()
	if n.Before(s.times[s.index].Add(s.frequency)) {
		return false
	}
	s.samples[s.index] = val
	s.times[s.index] = n
	s.index++
	if s.index == s.size {
		s.index = 0
	}
	return true
}

func (s FloatSampler) ExtrapolateValueLinearNano(time time.Time) float64 {
	prevI := s.index - 1
	if prevI < 0 {
		prevI += s.size
	}
	return s.samples[s.index] + float64(time.UnixNano())*(s.samples[s.index]-s.samples[prevI])/float64(s.times[s.index].UnixNano()-s.times[prevI].UnixNano())
}

func (s FloatSampler) ExtrapolateTimeLinearNano(val float64) time.Time {
	prevI := s.index - 1
	if prevI < 0 {
		prevI += s.size
	}
	t := float64(s.times[s.index].UnixNano()) + (val-s.samples[s.index])*float64(s.times[s.index].UnixNano()-s.times[prevI].UnixNano())/(s.samples[s.index]-s.samples[prevI])
	return time.Unix(0, int64(t))
}

func (s FloatSampler) ExtrapolateTimeLinearMillis(val float64) time.Time {
	prevI := s.index - 1
	if prevI < 0 {
		prevI += s.size
	}
	var slope float64 = 0
	for i := 1; i < s.size; i++ {
		index := (s.index + i) % s.size
		prevIndex := (s.index + i - 1) % s.size
		slope += float64(s.times[index].UnixMilli()-s.times[prevIndex].UnixMilli()) / (s.samples[index] - s.samples[prevIndex])
	}
	slope = slope / float64(s.size)
	t := float64(s.times[s.index].UnixMilli()) + (val-s.samples[s.index])*slope
	return time.UnixMilli(int64(t))
}

func (s FloatSampler) EstimateSecondsUntil(value float64) time.Duration {
	timeAtValue := s.ExtrapolateTimeLinearMillis(value)
	return timeAtValue.Sub(time.Now()).Round(time.Second)
}

func NewIntSampler(size int, frequency time.Duration) *IntSampler {
	return &IntSampler{
		FloatSampler: *NewFloatSampler(size, frequency),
	}
}

type IntSampler struct {
	FloatSampler
}

func (s *IntSampler) Add(value int) (added bool) {
	return s.FloatSampler.Add(float64(value))
}

func (s *IntSampler) EstimateSecondsUntil(val int) time.Duration {
	return s.FloatSampler.EstimateSecondsUntil(float64(val))
}
