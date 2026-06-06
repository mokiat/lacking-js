package internal

import (
	"github.com/mokiat/lacking/core/audio"
	"github.com/mokiat/wasmal"
)

type HighPassFilter struct {
	delegate wasmal.BiquadFilterNode
}

var _ audio.FrequencyFilter = (*HighPassFilter)(nil)
var _ Node = (*HighPassFilter)(nil)

func NewHighPassFilter(ctx wasmal.AudioContext) *HighPassFilter {
	delegate := ctx.CreateBiquadFilter()
	delegate.SetType(wasmal.BiquadFilterTypeHighpass)
	delegate.Frequency().SetValue(350.0)

	return &HighPassFilter{
		delegate: delegate,
	}
}

func (f *HighPassFilter) Input() wasmal.AudioNode {
	return f.delegate
}

func (f *HighPassFilter) Output() wasmal.AudioNode {
	return f.delegate
}

func (f *HighPassFilter) Frequency() float32 {
	return f.delegate.Frequency().Value()
}

func (f *HighPassFilter) SetFrequency(frequency float32) {
	f.delegate.Frequency().SetValue(frequency)
}
