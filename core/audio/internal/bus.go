package internal

import (
	"github.com/mokiat/gog/ds"
	"github.com/mokiat/lacking/core/audio"
	"github.com/mokiat/wasmal"
)

type Bus struct {
	gainFilter        *GainFilter
	compressionFilter *CompressionFilter
	reverbFilter      *ReverbFilter

	playbacks *ds.List[PlaybackNode]
	isPaused  bool
}

var _ audio.Bus = (*Bus)(nil)
var _ Node = (*Bus)(nil)

func NewBus(ctx wasmal.AudioContext, settings audio.BusSettings) *Bus {
	gainFilter := NewGainFilter(ctx)
	lastOutput := gainFilter.Output()

	var compressionFilter *CompressionFilter
	if settings.UseCompression {
		compressionFilter = NewCompressionFilter(ctx)
		lastOutput = compressionFilter.Output()
		gainFilter.Output().ConnectToNode(lastOutput)
	}

	var reverbFilter *ReverbFilter
	if settings.UseReverb {
		reverbFilter = NewReverbFilter(ctx)
		lastOutput.ConnectToNode(reverbFilter.Input())
	}

	return &Bus{
		gainFilter:        NewGainFilter(ctx),
		compressionFilter: compressionFilter,
		reverbFilter:      reverbFilter,

		playbacks: ds.EmptyList[PlaybackNode](),
	}
}

func (b *Bus) Input() wasmal.AudioNode {
	return b.gainFilter.Input()
}

func (b *Bus) Output() wasmal.AudioNode {
	if b.reverbFilter != nil {
		return b.reverbFilter.Output()
	}
	if b.compressionFilter != nil {
		return b.compressionFilter.Output()
	}
	return b.gainFilter.Output()
}

func (b *Bus) AddPlayback(p PlaybackNode) {
	if b.isPaused {
		p.InternalPause()
	}
	b.playbacks.Add(p)
	p.Output().ConnectToNode(b.Input())
}

func (b *Bus) RemovePlayback(p PlaybackNode) {
	b.playbacks.Remove(p)
	p.Output().Disconnect()
}

func (b *Bus) Gain() float32 {
	return b.gainFilter.Gain()
}

func (b *Bus) SetGain(gain float32) {
	b.gainFilter.SetGain(gain)
}

func (b *Bus) Compression() audio.Compression {
	if b.compressionFilter == nil {
		return nil
	}
	return b.compressionFilter
}

func (b *Bus) Reverb() audio.Reverb {
	if b.reverbFilter == nil {
		return nil
	}
	return b.reverbFilter
}

func (b *Bus) Pause() {
	if b.isPaused {
		return
	}
	b.isPaused = true
	for _, p := range b.playbacks.Unbox() {
		p.InternalPause()
	}
}

func (b *Bus) Resume() {
	if !b.isPaused {
		return
	}
	b.isPaused = false
	for _, p := range b.playbacks.Unbox() {
		p.InternalResume()
	}
}

func (b *Bus) Release() {
	b.Input().Disconnect()
	b.Output().Disconnect()
	b.playbacks = nil
}
