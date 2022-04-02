package internal

import (
	"errors"

	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/wasmgl"
)

func NewProgram(info render.ProgramInfo) *Program {
	program := &Program{
		raw: wasmgl.CreateProgram(),
	}
	if vertexShader, ok := info.VertexShader.(*Shader); ok {
		wasmgl.AttachShader(program.raw, vertexShader.raw)
		defer wasmgl.DetachShader(program.raw, vertexShader.raw)
	}
	if fragmentShader, ok := info.FragmentShader.(*Shader); ok {
		wasmgl.AttachShader(program.raw, fragmentShader.raw)
		defer wasmgl.DetachShader(program.raw, fragmentShader.raw)
	}
	if err := program.link(); err != nil {
		log.Error("Program link error: %v", err)
	}
	return program
}

type Program struct {
	raw wasmgl.Program
}

func (p *Program) UniformLocation(name string) render.UniformLocation {
	return wasmgl.GetUniformLocation(p.raw, name)
}

func (p *Program) Release() {
	wasmgl.DeleteProgram(p.raw)
	p.raw = wasmgl.NilProgram
}

func (p *Program) link() error {
	wasmgl.LinkProgram(p.raw)
	if !p.isLinkSuccessful() {
		return errors.New(p.getInfoLog())
	}
	return nil
}

func (p *Program) isLinkSuccessful() bool {
	result := wasmgl.GetProgramParameter(p.raw, wasmgl.LINK_STATUS)
	return result.Bool()
}

func (p *Program) getInfoLog() string {
	return wasmgl.GetProgramInfoLog(p.raw)
}
