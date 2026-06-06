package internal

import (
	"github.com/mokiat/lacking/core/audio"
	"github.com/mokiat/wasmal"
)

type ReverbFilter struct {

	// FIXME: Using GainNode as placeholder.
	delegate wasmal.GainNode
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
	panic("TODO")
}

func (f *ReverbFilter) SetRoomSize(size float32) {
	panic("TODO")
}

func (f *ReverbFilter) Damping() float32 {
	panic("TODO")
}

func (f *ReverbFilter) SetDamping(damping float32) {
	panic("TODO")
}

func (f *ReverbFilter) Dry() float32 {
	panic("TODO")
}

func (f *ReverbFilter) SetDry(dry float32) {
	panic("TODO")
}

func (f *ReverbFilter) Wet() float32 {
	panic("TODO")
}

func (f *ReverbFilter) SetWet(wet float32) {
	panic("TODO")
}
