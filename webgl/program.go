package webgl

import (
	"fmt"

	"github.com/mokiat/wasmgl"
)

func NewProgram() *Program {
	return &Program{}
}

type Program struct {
	raw wasmgl.Program
}

func (p *Program) Raw() wasmgl.Program {
	return p.raw
}

func (p *Program) Allocate(info ProgramAllocateInfo) {
	p.raw = wasmgl.CreateProgram()
	if info.VertexShader != nil {
		wasmgl.AttachShader(p.raw, info.VertexShader.Raw())
		defer wasmgl.DetachShader(p.raw, info.VertexShader.Raw())
	}
	if info.FragmentShader != nil {
		wasmgl.AttachShader(p.raw, info.FragmentShader.Raw())
		defer wasmgl.DetachShader(p.raw, info.FragmentShader.Raw())
	}
	wasmgl.LinkProgram(p.raw)
	if !p.isLinkSuccessful() {
		panic(fmt.Errorf("failed to link program: %s", p.getInfoLog()))
	}
}

func (p *Program) UniformLocation(name string) wasmgl.UniformLocation {
	return wasmgl.GetUniformLocation(p.raw, name)
}

func (p *Program) Use() {
	wasmgl.UseProgram(p.raw)
}

func (p *Program) Release() {
	wasmgl.DeleteProgram(p.raw)
	p.raw = wasmgl.Program{}
}

func (p *Program) isLinkSuccessful() bool {
	result := wasmgl.GetProgramParameter(p.raw, wasmgl.LINK_STATUS)
	return result.Bool()
}

func (p *Program) getInfoLog() string {
	return wasmgl.GetProgramInfoLog(p.raw)
}

type ProgramAllocateInfo struct {
	VertexShader                 *Shader
	TessellationControlShader    *Shader
	TessellationEvaluationShader *Shader
	FragmentShader               *Shader
}
