package internal

import (
	"github.com/mokiat/lacking/audio"
	"github.com/mokiat/wasmal"
)

type CompressorNode struct {
	audio.Node // marker interface
	delegate   wasmal.DynamicsCompressorNode
}

var _ audio.CompressorNode = (*CompressorNode)(nil)
var _ Node = (*CompressorNode)(nil)

func newCompressorNode(audioContext wasmal.AudioContext) *CompressorNode {
	return &CompressorNode{
		delegate: audioContext.CreateDynamicsCompressor(),
	}
}

func (n *CompressorNode) AudioNode() wasmal.AudioNode {
	return n.delegate
}

func (n *CompressorNode) Attack() float32 {
	return float32(n.delegate.Attack().Value())
}

func (n *CompressorNode) SetAttack(attack float32) {
	n.delegate.Attack().SetValue(float64(attack))
}

func (n *CompressorNode) Release() float32 {
	return float32(n.delegate.Release().Value())
}

func (n *CompressorNode) SetRelease(release float32) {
	n.delegate.Release().SetValue(float64(release))
}

func (n *CompressorNode) Ratio() float32 {
	return float32(n.delegate.Ratio().Value())
}

func (n *CompressorNode) SetRatio(ratio float32) {
	n.delegate.Ratio().SetValue(float64(ratio))
}

func (n *CompressorNode) Knee() float32 {
	return float32(n.delegate.Knee().Value())
}

func (n *CompressorNode) SetKnee(knee float32) {
	n.delegate.Knee().SetValue(float64(knee))
}

func (n *CompressorNode) Threshold() float32 {
	return float32(n.delegate.Threshold().Value())
}

func (n *CompressorNode) SetThreshold(threshold float32) {
	n.delegate.Threshold().SetValue(float64(threshold))
}

func (n *CompressorNode) Delete() {
	if n.delegate != nil {
		n.delegate.Disconnect()
		n.delegate = nil
	}
}
