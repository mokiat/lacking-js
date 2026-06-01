package internal

import (
	"github.com/mokiat/lacking/audio"
	"github.com/mokiat/wasmal"
)

type PanNode struct {
	audio.Node // marker interface
	delegate   wasmal.StereoPannerNode
}

var _ audio.PanNode = (*PanNode)(nil)
var _ Node = (*PanNode)(nil)

func newPanNode(audioContext wasmal.AudioContext) *PanNode {
	return &PanNode{
		delegate: audioContext.CreateStereoPanner(),
	}
}

func (n *PanNode) AudioNode() wasmal.AudioNode {
	return n.delegate
}

func (n *PanNode) Pan() float32 {
	return float32(n.delegate.Pan().Value())
}

func (n *PanNode) SetPan(pan float32) {
	n.delegate.Pan().SetValue(float64(pan))
}

func (n *PanNode) Delete() {
	if n.delegate != nil {
		n.delegate.Disconnect()
		n.delegate = nil
	}
}
