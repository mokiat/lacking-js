package internal

import "github.com/mokiat/lacking/audio"

type DefaultPlayback struct {
	*BasePlayback
}

var _ audio.Playback = (*DefaultPlayback)(nil)

func NewDefaultPlayback(basePlayback *BasePlayback) *DefaultPlayback {
	return &DefaultPlayback{
		BasePlayback: basePlayback,
	}
}

func (p *DefaultPlayback) Release() {
	defer p.BasePlayback.Release()
	p.Output().Disconnect()
}
