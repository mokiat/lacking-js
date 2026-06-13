package internal

import (
	"github.com/mokiat/lacking/core/audio"
	"github.com/mokiat/wasmal"
)

type MasterBus struct {
	gainFilter        *GainFilter
	compressionFilter *CompressionFilter
}

var _ audio.MasterBus = (*MasterBus)(nil)
var _ Node = (*MasterBus)(nil)

func NewMasterBus(ctx wasmal.AudioContext) *MasterBus {
	gainFilter := NewGainFilter(ctx)
	compressionFilter := NewCompressionFilter(ctx)

	gainFilter.Output().ConnectToNode(compressionFilter.Input())

	return &MasterBus{
		gainFilter:        gainFilter,
		compressionFilter: compressionFilter,
	}
}

func (b *MasterBus) Input() wasmal.AudioNode {
	return b.gainFilter.Input()
}

func (b *MasterBus) Output() wasmal.AudioNode {
	return b.compressionFilter.Output()
}

func (b *MasterBus) Gain() float32 {
	return b.gainFilter.Gain()
}

func (b *MasterBus) SetGain(gain float32) {
	b.gainFilter.SetGain(gain)
}

func (b *MasterBus) Compression() audio.Compression {
	return b.compressionFilter
}
