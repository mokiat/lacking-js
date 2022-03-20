package webgl

import (
	"fmt"
	"strings"

	"github.com/mokiat/wasmgl"
)

func NewShader() *Shader {
	return &Shader{}
}

type Shader struct {
	raw wasmgl.Shader
}

func (s *Shader) Raw() wasmgl.Shader {
	return s.raw
}

func (s *Shader) Allocate(info ShaderAllocateInfo) {
	s.raw = wasmgl.CreateShader(info.ShaderType)
	wasmgl.ShaderSource(s.raw, info.SourceCode)
	wasmgl.CompileShader(s.raw)
	if !s.isCompileSuccessful() {
		panic(fmt.Errorf("failed to compile shader: %s", s.getInfoLog()))
	}
}

func (s *Shader) Release() {
	wasmgl.DeleteShader(s.raw)
	s.raw = wasmgl.Shader{}
}

func (s *Shader) isCompileSuccessful() bool {
	result := wasmgl.GetShaderParameter(s.raw, wasmgl.COMPILE_STATUS)
	return result.Bool()
}

func (s *Shader) getInfoLog() string {
	return wasmgl.GetShaderInfoLog(s.raw)
}

type ShaderAllocateInfo struct {
	ShaderType int
	SourceCode string
}

func NewShaderSourceBuilder(template string) *ShaderSourceBuilder {
	return &ShaderSourceBuilder{
		version:        "300 es",
		floatPrecision: "mediump",
		features:       []string{},
		template:       template,
	}
}

type ShaderSourceBuilder struct {
	version        string
	floatPrecision string
	features       []string
	template       string
}

func (b *ShaderSourceBuilder) SetFloatPrecision(precision string) {
	b.floatPrecision = precision
}

func (b *ShaderSourceBuilder) SetVersion(version string) {
	b.version = version
}

func (b *ShaderSourceBuilder) AddFeature(feature string) {
	b.features = append(b.features, feature)
}

func (b *ShaderSourceBuilder) Build() string {
	var builder strings.Builder
	builder.WriteString("#version ")
	builder.WriteString(b.version)
	builder.WriteRune('\n')
	builder.WriteString("precision ")
	builder.WriteString(b.floatPrecision)
	builder.WriteString(" float;")
	builder.WriteRune('\n')
	for _, feature := range b.features {
		builder.WriteString("#define ")
		builder.WriteString(feature)
		builder.WriteRune('\n')
	}
	builder.WriteString(b.template)
	builder.WriteRune('\n')
	return builder.String()
}
