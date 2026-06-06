package internal

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/audio"
	"github.com/mokiat/wasmal"
)

type SpatialNode struct {
	audio.Node // marker interface
	delegate   wasmal.PannerNode
}

var _ audio.SpatialNode = (*SpatialNode)(nil)
var _ Node = (*SpatialNode)(nil)

func newSpatialNode(audioContext wasmal.AudioContext) *SpatialNode {
	return &SpatialNode{
		delegate: audioContext.CreatePanner(),
	}
}

func (n *SpatialNode) AudioNode() wasmal.AudioNode {
	return n.delegate
}

func (n *SpatialNode) Position() sprec.Vec3 {
	return sprec.Vec3{
		X: float32(n.delegate.PositionX().Value()),
		Y: float32(n.delegate.PositionY().Value()),
		Z: float32(n.delegate.PositionZ().Value()),
	}
}

func (n *SpatialNode) SetPosition(position sprec.Vec3) {
	n.delegate.PositionX().SetValue(position.X)
	n.delegate.PositionY().SetValue(position.Y)
	n.delegate.PositionZ().SetValue(position.Z)
}

func (n *SpatialNode) Delete() {
	if n.delegate != nil {
		n.delegate.Disconnect()
		n.delegate = nil
	}
}
