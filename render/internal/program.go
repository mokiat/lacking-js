package internal

import (
	"errors"

	"github.com/mokiat/lacking/render"
	"github.com/mokiat/wasmgl"
)

type ProgramInfo struct {
	Label           string
	VertexCode      string
	FragmentCode    string
	TextureBindings []render.TextureBinding
	UniformBindings []render.UniformBinding
}

func NewProgram(info ProgramInfo) *Program {
	vertexShader := newVertexShader(info.Label, info.VertexCode)
	defer vertexShader.Release()

	fragmentShader := newFragmentShader(info.Label, info.FragmentCode)
	defer fragmentShader.Release()

	program := &Program{
		label: info.Label,
		raw:   wasmgl.CreateProgram(),
	}

	wasmgl.AttachShader(program.raw, vertexShader.raw)
	defer wasmgl.DetachShader(program.raw, vertexShader.raw)

	wasmgl.AttachShader(program.raw, fragmentShader.raw)
	defer wasmgl.DetachShader(program.raw, fragmentShader.raw)

	if err := program.link(); err != nil {
		logger.Error("Program (%v) link error: %v", info.Label, err)
	}

	if len(info.TextureBindings) > 0 {
		wasmgl.UseProgram(program.raw)
		for _, binding := range info.TextureBindings {
			location := wasmgl.GetUniformLocation(program.raw, binding.Name)
			if location.IsValid() {
				wasmgl.Uniform1i(location, wasmgl.GLint(binding.Index))
			}
		}
		wasmgl.UseProgram(wasmgl.NilProgram)
	}

	for _, binding := range info.UniformBindings {
		location := wasmgl.GetUniformBlockIndex(program.raw, binding.Name)
		if location != wasmgl.INVALID_INDEX {
			wasmgl.UniformBlockBinding(program.raw, location, wasmgl.GLuint(binding.Index))
		}
	}

	program.id = programs.Allocate(program)
	return program
}

type Program struct {
	render.ProgramMarker

	label string
	id    uint32
	raw   wasmgl.Program
}

func (p *Program) Label() string {
	return p.label
}

func (p *Program) Release() {
	programs.Release(p.id)
	wasmgl.DeleteProgram(p.raw)
	p.raw = wasmgl.NilProgram
	p.id = 0
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
	return result.GLboolean()
}

func (p *Program) getInfoLog() string {
	return wasmgl.GetProgramInfoLog(p.raw)
}
