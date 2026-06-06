package internal

import (
	"github.com/mokiat/lacking/audio"
	"github.com/mokiat/wasmal"
)

type DelayNode struct {
	audio.Node // marker interface
	delegate   wasmal.DelayNode
}

var _ audio.DelayNode = (*DelayNode)(nil)
var _ Node = (*DelayNode)(nil)

func newDelayNode(audioContext wasmal.AudioContext) *DelayNode {
	return &DelayNode{
		delegate: audioContext.CreateDelay(1.0),
	}
}

func (n *DelayNode) AudioNode() wasmal.AudioNode {
	return n.delegate
}

func (n *DelayNode) DelayTime() float32 {
	return float32(n.delegate.DelayTime().Value())
}

func (n *DelayNode) SetDelayTime(delayTime float32) {
	n.delegate.DelayTime().SetValue(delayTime)
}

func (n *DelayNode) Delete() {
	if n.delegate != nil {
		n.delegate.Disconnect()
		n.delegate = nil
	}
}
