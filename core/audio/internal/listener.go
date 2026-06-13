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
	l.delegate.PositionX().SetValue(position.X)
	l.delegate.PositionY().SetValue(position.Y)
	l.delegate.PositionZ().SetValue(position.Z)
}

func (l *SpatialListener) Rotation() sprec.Quat {
	return l.rotation
}

func (l *SpatialListener) SetRotation(rotation sprec.Quat) {
	l.rotation = rotation

	orientationY := rotation.OrientationY()
	l.delegate.UpX().SetValue(orientationY.X)
	l.delegate.UpY().SetValue(orientationY.Y)
	l.delegate.UpZ().SetValue(orientationY.Z)

	orientationZ := rotation.OrientationZ()
	l.delegate.ForwardX().SetValue(orientationZ.X)
	l.delegate.ForwardY().SetValue(orientationZ.Y)
	l.delegate.ForwardZ().SetValue(orientationZ.Z)
}
