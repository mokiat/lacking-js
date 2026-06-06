package internal

import (
	"github.com/mokiat/lacking/core/audio"
	"github.com/mokiat/wasmal"
)

type ReverbFilter struct {
	// FIXME: Using GainNode as placeholder.
	delegate wasmal.GainNode

	roomSize float32
	damping  float32
	dry      float32
	wet      float32
}

var _ audio.Reverb = (*ReverbFilter)(nil)
var _ Node = (*ReverbFilter)(nil)

func NewReverbFilter(ctx wasmal.AudioContext) *ReverbFilter {
	return &ReverbFilter{
		delegate: ctx.CreateGain(),
	}
}

func (f *ReverbFilter) Input() wasmal.AudioNode {
	return f.delegate
}

func (f *ReverbFilter) Output() wasmal.AudioNode {
	return f.delegate
}

func (f *ReverbFilter) RoomSize() float32 {
	return f.roomSize
}

func (f *ReverbFilter) SetRoomSize(size float32) {
	f.roomSize = size
}

func (f *ReverbFilter) Damping() float32 {
	return f.damping
}

func (f *ReverbFilter) SetDamping(damping float32) {
	f.damping = damping
}

func (f *ReverbFilter) Dry() float32 {
	return f.dry
}

func (f *ReverbFilter) SetDry(dry float32) {
	f.dry = dry
}

func (f *ReverbFilter) Wet() float32 {
	return f.wet
}

func (f *ReverbFilter) SetWet(wet float32) {
	f.wet = wet
}
