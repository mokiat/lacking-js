package internal

import (
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/wasmgl"
)

func NewQueue() *Queue {
	return &Queue{}
}

type Queue struct {
}

func (q *Queue) Invalidate() {
	// TODO
}

func (q *Queue) WriteBuffer(buffer render.Buffer, offset int, data []byte) {
	actualBuffer := buffer.(*Buffer)
	wasmgl.BindBuffer(actualBuffer.kind, actualBuffer.raw)
	wasmgl.BufferSubData(actualBuffer.kind, wasmgl.GLintptr(offset), data)
	wasmgl.BindBuffer(actualBuffer.kind, wasmgl.NilBuffer)
}
