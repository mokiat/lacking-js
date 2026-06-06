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

	channelData := buffer.GetChannelData(0)
	for i, frame := range frames {
		channelData.Set(i, frame.Left)
	}

	channelData = buffer.GetChannelData(1)
	for i, frame := range frames {
		channelData.Set(i, frame.Right)
	}

	return &Media{
		buffer: buffer,
		length: buffer.Duration(),
	}
}

func (m *Media) Length() float64 {
	return m.length
}

func (m *Media) Release() {}
