package internal

import (
	"github.com/mokiat/lacking/core/audio"
	"github.com/mokiat/wasmal"
)

type LowPassFilter struct {
	delegate wasmal.BiquadFilterNode
}

var _ audio.FrequencyFilter = (*LowPassFilter)(nil)
var _ Node = (*LowPassFilter)(nil)

func NewLowPassFilter(ctx wasmal.AudioContext) *LowPassFilter {
	delegate := ctx.CreateBiquadFilter()
	delegate.SetType(wasmal.BiquadFilterTypeLowpass)
	delegate.Frequency().SetValue(350.0)

	return &LowPassFilter{
		delegate: delegate,
	}
}

func (f *LowPassFilter) Input() wasmal.AudioNode {
	return f.delegate
}

func (f *LowPassFilter) Output() wasmal.AudioNode {
	return f.delegate
}

func (f *LowPassFilter) Frequency() float32 {
	return f.delegate.Frequency().Value()
}

func (f *LowPassFilter) SetFrequency(frequency float32) {
	f.delegate.Frequency().SetValue(frequency)
}
