package internal

import (
	"github.com/mokiat/lacking/core/audio"
	"github.com/mokiat/wasmal"
)

type CompressionFilter struct {
	delegate wasmal.DynamicsCompressorNode
}

var _ audio.Compression = (*CompressionFilter)(nil)
var _ Node = (*CompressionFilter)(nil)

func NewCompressionFilter(ctx wasmal.AudioContext) *CompressionFilter {
	return &CompressionFilter{
		delegate: ctx.CreateDynamicsCompressor(),
	}
}

func (c *CompressionFilter) Input() wasmal.AudioNode {
	return c.delegate
}

func (c *CompressionFilter) Output() wasmal.AudioNode {
	return c.delegate
}

func (c *CompressionFilter) Attack() float32 {
	return float32(c.delegate.Attack().Value())
}

func (c *CompressionFilter) SetAttack(attack float32) {
	c.delegate.Attack().SetValue(attack)
}

func (c *CompressionFilter) Release() float32 {
	return float32(c.delegate.Release().Value())
}

func (c *CompressionFilter) SetRelease(release float32) {
	c.delegate.Release().SetValue(release)
}

func (c *CompressionFilter) Ratio() float32 {
	return float32(c.delegate.Ratio().Value())
}

func (c *CompressionFilter) SetRatio(ratio float32) {
	c.delegate.Ratio().SetValue(ratio)
}

func (c *CompressionFilter) Knee() float32 {
	return float32(c.delegate.Knee().Value())
}

func (c *CompressionFilter) SetKnee(knee float32) {
	c.delegate.Knee().SetValue(knee)
}

func (c *CompressionFilter) Threshold() float32 {
	return float32(c.delegate.Threshold().Value())
}

func (c *CompressionFilter) SetThreshold(threshold float32) {
	c.delegate.Threshold().SetValue(threshold)
}
