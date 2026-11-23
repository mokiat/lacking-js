package internal

import (
	"github.com/mokiat/lacking/core/audio"
	"github.com/mokiat/wasmal"
)

type Media struct {
	buffer wasmal.AudioBuffer
	length float64
}

var _ audio.Media = (*Media)(nil)

func NewMedia(ctx wasmal.AudioContext, data audio.MediaData) *Media {
	frames := data.Frames

	buffer := ctx.CreateBuffer(2, uint32(len(frames)), float32(data.SampleRate))

	left := make([]float32, len(frames))
	right := make([]float32, len(frames))
	for i, frame := range frames {
		left[i] = frame.Left
		right[i] = frame.Right
	}
	buffer.CopyToChannel(left, 0, 0)
	buffer.CopyToChannel(right, 1, 0)

	return &Media{
		buffer: buffer,
		length: buffer.Duration(),
	}
}

func (m *Media) Length() float64 {
	return m.length
}

func (m *Media) Release() {}
