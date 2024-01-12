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
		raw:      wasmgl.CreateProgram(),
		uniforms: make(map[*UniformLocation]struct{}),
	}

	wasmgl.AttachShader(program.raw, vertexShader.raw)
	defer wasmgl.DetachShader(program.raw, vertexShader.raw)

	wasmgl.AttachShader(program.raw, fragmentShader.raw)
	defer wasmgl.DetachShader(program.raw, fragmentShader.raw)

	if err := program.link(); err != nil {
		logger.Error("Program link error: %v!", err)
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
	render.ProgramObject
	id       uint32
	raw      wasmgl.Program
	uniforms map[*UniformLocation]struct{}
}

// Deprecated: To be removed.
func (p *Program) UniformLocation(name string) render.UniformLocation {
	result := &UniformLocation{
		raw: wasmgl.GetUniformLocation(p.raw, name),
	}
	result.id = int32(locations.Allocate(result))
	p.uniforms[result] = struct{}{}
	return result
}

func (p *Program) Release() {
	for uniform := range p.uniforms {
		locations.Release(uint32(uniform.id))
	}
	p.uniforms = nil
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

type UniformLocation struct {
	id  int32
	raw wasmgl.UniformLocation
}
