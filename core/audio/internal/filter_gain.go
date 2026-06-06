package internal

import "github.com/mokiat/wasmal"

type GainFilter struct {
	delegate wasmal.GainNode
}

var _ Node = (*GainFilter)(nil)

func NewGainFilter(ctx wasmal.AudioContext) *GainFilter {
	return &GainFilter{
		delegate: ctx.CreateGain(),
	}
}

func (f *GainFilter) Input() wasmal.AudioNode {
	return f.delegate
}

func (f *GainFilter) Output() wasmal.AudioNode {
	return f.delegate
}

func (f *GainFilter) Gain() float32 {
	return float32(f.delegate.Gain().Value())
}

func (f *GainFilter) SetGain(gain float32) {
	f.delegate.Gain().SetValue(gain)
}
