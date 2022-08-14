package internal

import "strings"

func NewShaderSourceBuilder(template string) *ShaderSourceBuilder {
	return &ShaderSourceBuilder{
		version:                  "300 es",
		floatPrecision:           "mediump",
		sampler2DShadowPrecision: "mediump",
		features:                 []string{},
		template:                 template,
	}
}

type ShaderSourceBuilder struct {
	version                  string
	floatPrecision           string
	sampler2DShadowPrecision string
	features                 []string
	template                 string
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
	builder.WriteString("precision ")
	builder.WriteString(b.sampler2DShadowPrecision)
	builder.WriteString(" sampler2DShadow;")
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
