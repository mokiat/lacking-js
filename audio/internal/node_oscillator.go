package internal

import (
	"github.com/mokiat/lacking/audio"
	"github.com/mokiat/wasmal"
)

type OscillatorNode struct {
	audio.Node // marker interface
	delegate   wasmal.OscillatorNode
}

var _ audio.OscillatorNode = (*OscillatorNode)(nil)
var _ Node = (*OscillatorNode)(nil)

func newOscillatorNode(audioContext wasmal.AudioContext) *OscillatorNode {
	return &OscillatorNode{
		delegate: audioContext.CreateOscillator(),
	}
}

func (n *OscillatorNode) AudioNode() wasmal.AudioNode {
	return n.delegate
}

func (n *OscillatorNode) Frequency() float32 {
	return float32(n.delegate.Frequency().Value())
}

func (n *OscillatorNode) SetFrequency(frequency float32) {
	n.delegate.Frequency().SetValue(frequency)
}

func (n *OscillatorNode) Delete() {
	if n.delegate != nil {
		n.delegate.Disconnect()
		n.delegate = nil
	}
}
