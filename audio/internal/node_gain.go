package internal

import (
	"github.com/mokiat/lacking/audio"
	"github.com/mokiat/wasmal"
)

type GainNode struct {
	audio.Node // marker interface
	delegate   wasmal.GainNode
}

var _ audio.GainNode = (*GainNode)(nil)
var _ Node = (*GainNode)(nil)

func newGainNode(audioContext wasmal.AudioContext) *GainNode {
	return &GainNode{
		delegate: audioContext.CreateGain(),
	}
}

func (n *GainNode) AudioNode() wasmal.AudioNode {
	return n.delegate
}

func (n *GainNode) Gain() float32 {
	return float32(n.delegate.Gain().Value())
}

func (n *GainNode) SetGain(gain float32) {
	n.delegate.Gain().SetValue(float64(gain))
}

func (n *GainNode) Delete() {
	if n.delegate != nil {
		n.delegate.Disconnect()
		n.delegate = nil
	}
}
