package internal

import (
	"github.com/mokiat/lacking/audio"
	"github.com/mokiat/wasmal"
)

type OutputNode struct {
	audio.Node  // marker interface
	destination wasmal.AudioDestinationNode
}

var _ audio.Node = (*OutputNode)(nil)
var _ Node = (*OutputNode)(nil)

func (n *OutputNode) AudioNode() wasmal.AudioNode {
	return n.destination
}

func newOutputNode(audioContext wasmal.AudioContext) *OutputNode {
	return &OutputNode{
		destination: audioContext.Destination(),
	}
}
