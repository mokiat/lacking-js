package internal

import (
	"github.com/mokiat/lacking/audio"
	"github.com/mokiat/wasmal"
)

type ConnectorNode struct {
	audio.Node // marker interface
	delegate   wasmal.AudioNode
}

var _ audio.ConnectorNode = (*ConnectorNode)(nil)
var _ Node = (*ConnectorNode)(nil)

func newConnectorNode(audioContext wasmal.AudioContext) *ConnectorNode {
	return &ConnectorNode{
		delegate: audioContext.CreateGain(),
	}
}

func (n *ConnectorNode) AudioNode() wasmal.AudioNode {
	return n.delegate
}

func (n *ConnectorNode) Delete() {
	if n.delegate != nil {
		n.delegate.Disconnect()
		n.delegate = nil
	}
}
