package internal

import (
	"errors"

	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/wasmgl"
)

func NewProgram(info render.ProgramInfo) *Program {
	program := &Program{
		raw:      wasmgl.CreateProgram(),
		uniforms: make(map[*UniformLocation]struct{}),
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
	program.id = programs.Allocate(program)
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
	return program
}

type Program struct {
	render.ProgramObject
	id       uint32
	raw      wasmgl.Program
	uniforms map[*UniformLocation]struct{}
}

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
