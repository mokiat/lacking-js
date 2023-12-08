package internal

import "github.com/mokiat/wasmal"

type Playback struct {
	node wasmal.AudioScheduledSourceNode
}

func (p *Playback) Stop() {
	if p.node != nil {
		p.node.Stop(0.0)
		p.node = nil
	}
}
