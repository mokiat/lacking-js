package internal

import "github.com/mokiat/wasmgl"

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
