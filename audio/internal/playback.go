package internal

type Playback struct {
	srcNode  *PlaybackNode
	panNode  *PanNode
	gainNode *GainNode
}

func (p *Playback) Stop() {
	if p.srcNode != nil {
		p.srcNode.Stop()
	}
}
