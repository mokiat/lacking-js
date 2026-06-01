package internal

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/audio"
	"github.com/mokiat/wasmal"
)

type SpatialListener struct {
	delegate wasmal.AudioListener
}

var _ audio.SpatialListener = (*SpatialListener)(nil)

func newSpatialListener(audioContext wasmal.AudioContext) *SpatialListener {
	return &SpatialListener{
		delegate: audioContext.Listener(),
	}
}

func (l *SpatialListener) Position() sprec.Vec3 {
	return sprec.Vec3{
		X: float32(l.delegate.PositionX().Value()),
		Y: float32(l.delegate.PositionY().Value()),
		Z: float32(l.delegate.PositionZ().Value()),
	}
}

func (l *SpatialListener) SetPosition(position sprec.Vec3) {
	l.delegate.PositionX().SetValue(float64(position.X))
	l.delegate.PositionY().SetValue(float64(position.Y))
	l.delegate.PositionZ().SetValue(float64(position.Z))
}

func (l *SpatialListener) Rotation() sprec.Quat {
	orientationX := sprec.Vec3{
		X: float32(l.delegate.ForwardX().Value()),
		Y: float32(l.delegate.ForwardY().Value()),
		Z: float32(l.delegate.ForwardZ().Value()),
	}
	orientationY := sprec.Vec3{
		X: float32(l.delegate.UpX().Value()),
		Y: float32(l.delegate.UpY().Value()),
		Z: float32(l.delegate.UpZ().Value()),
	}
	orientationZ := sprec.Vec3Cross(orientationX, orientationY)

	mat := sprec.OrientationMat4(orientationX, orientationY, orientationZ)
	return mat.Rotation()
}

func (l *SpatialListener) SetRotation(rotation sprec.Quat) {
	orientationX := rotation.OrientationX()
	l.delegate.ForwardX().SetValue(float64(orientationX.X))
	l.delegate.ForwardY().SetValue(float64(orientationX.Y))
	l.delegate.ForwardZ().SetValue(float64(orientationX.Z))

	orientationY := rotation.OrientationY()
	l.delegate.UpX().SetValue(float64(orientationY.X))
	l.delegate.UpY().SetValue(float64(orientationY.Y))
	l.delegate.UpZ().SetValue(float64(orientationY.Z))
}
