package internal

import (
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/wasmgl"
)

func NewLimits() *Limits {
	uniformBufferOffsetAlignment := wasmgl.GetParameter(wasmgl.UNIFORM_BUFFER_OFFSET_ALIGNMENT).GLint()
	return &Limits{
		uniformBufferOffsetAlignment: int(uniformBufferOffsetAlignment),
	}
}

type Limits struct {
	uniformBufferOffsetAlignment int
}

func (l Limits) UniformBufferOffsetAlignment() int {
	return l.uniformBufferOffsetAlignment
}

func (l Limits) Quality() render.Quality {
	return render.QualityHigh
}
