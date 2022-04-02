package internal

import "github.com/mokiat/lacking/render"

func NewPipeline(info render.PipelineInfo) *Pipeline {
	return &Pipeline{
		PipelineInfo: info,
	}
}

type Pipeline struct {
	render.PipelineInfo
}

func (p *Pipeline) Release() {

}
