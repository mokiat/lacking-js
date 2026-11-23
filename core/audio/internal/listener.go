package internal

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/core/audio"
	"github.com/mokiat/wasmal"
)

type SpatialListener struct {
	delegate wasmal.AudioListener

	position sprec.Vec3
	rotation sprec.Quat
}

var _ audio.SpatialListener = (*SpatialListener)(nil)

func NewSpatialListener(ctx wasmal.AudioContext) *SpatialListener {
	return &SpatialListener{
		delegate: ctx.Listener(),

		position: sprec.ZeroVec3(),
		rotation: sprec.IdentityQuat(),
	}
}

func (l *SpatialListener) Position() sprec.Vec3 {
	return l.position
}

func (l *SpatialListener) SetPosition(position sprec.Vec3) {
	l.position = position
	l.delegate.SetPosition(position.X, position.Y, position.Z)
}

func (l *SpatialListener) Rotation() sprec.Quat {
	return l.rotation
}

func (l *SpatialListener) SetRotation(rotation sprec.Quat) {
	l.rotation = rotation

	up := rotation.OrientationY()
	// The listener looks down its -Z axis.
	forward := sprec.InverseVec3(rotation.OrientationZ())
	l.delegate.SetOrientation(
		forward.X, forward.Y, forward.Z,
		up.X, up.Y, up.Z,
	)
}
