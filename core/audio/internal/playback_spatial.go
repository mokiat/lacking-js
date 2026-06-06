package internal

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/core/audio"
	"github.com/mokiat/wasmal"
)

type SpatialPlayback struct {
	*BasePlayback
	pannerNode wasmal.PannerNode

	position sprec.Vec3
	rotation sprec.Quat
}

var _ audio.SpatialPlayback = (*SpatialPlayback)(nil)

func NewSpatialPlayback(basePlayback *BasePlayback) *SpatialPlayback {
	ctx := basePlayback.ctx
	pannerNode := ctx.CreatePanner()
	basePlayback.Output().ConnectToNode(pannerNode)

	return &SpatialPlayback{
		BasePlayback: basePlayback,
		pannerNode:   pannerNode,

		position: sprec.ZeroVec3(),
		rotation: sprec.IdentityQuat(),
	}
}

func (p *SpatialPlayback) Output() wasmal.AudioNode {
	return p.pannerNode
}

// Release removes this playback from its bus, triggering a pipeline rebuild.
func (p *SpatialPlayback) Release() {
	defer p.BasePlayback.Release()
	p.Output().Disconnect()
}

// Position returns the emitter's position in world space.
func (p *SpatialPlayback) Position() sprec.Vec3 {
	return p.position
}

// SetPosition sets the emitter's position in world space.
func (p *SpatialPlayback) SetPosition(position sprec.Vec3) {
	p.position = position
	p.pannerNode.PositionX().SetValue(position.X)
	p.pannerNode.PositionY().SetValue(position.Y)
	p.pannerNode.PositionZ().SetValue(position.Z)
}

// Rotation returns the emitter's orientation in world space. The emitter's
// forward (cone) direction is its +Z axis.
func (p *SpatialPlayback) Rotation() sprec.Quat {
	return p.rotation
}

// SetRotation sets the emitter's orientation in world space.
func (p *SpatialPlayback) SetRotation(rotation sprec.Quat) {
	p.rotation = rotation
	orientationZ := rotation.OrientationZ()
	p.pannerNode.OrientationX().SetValue(orientationZ.X)
	p.pannerNode.OrientationY().SetValue(orientationZ.Y)
	p.pannerNode.OrientationZ().SetValue(orientationZ.Z)
}

// InnerConeAngle returns the inner cone angle of the emitter.
func (p *SpatialPlayback) InnerConeAngle() sprec.Angle {
	return sprec.Degrees(float32(p.pannerNode.ConeInnerAngle()))
}

// SetInnerConeAngle sets the inner cone angle of the emitter.
func (p *SpatialPlayback) SetInnerConeAngle(angle sprec.Angle) {
	p.pannerNode.SetConeInnerAngle(float64(angle.Degrees()))
}

// OuterConeAngle returns the outer cone angle of the emitter.
func (p *SpatialPlayback) OuterConeAngle() sprec.Angle {
	return sprec.Degrees(float32(p.pannerNode.ConeOuterAngle()))
}

// SetOuterConeAngle sets the outer cone angle of the emitter.
func (p *SpatialPlayback) SetOuterConeAngle(angle sprec.Angle) {
	p.pannerNode.SetConeOuterAngle(float64(angle.Degrees()))
}

// OuterConeGain returns the gain applied when the listener is outside the
// outer cone.
func (p *SpatialPlayback) OuterConeGain() float32 {
	return float32(p.pannerNode.ConeOuterGain())
}

// SetOuterConeGain sets the gain applied when the listener is outside the
// outer cone.
func (p *SpatialPlayback) SetOuterConeGain(gain float32) {
	p.pannerNode.SetConeOuterGain(float64(gain))
}
