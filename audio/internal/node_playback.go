package internal

import (
	"github.com/mokiat/lacking/audio"
	"github.com/mokiat/wasmal"
)

type PlaybackNode struct {
	audio.Node   // marker interface
	audioContext wasmal.AudioContext
	buffer       wasmal.AudioBuffer
	delegate     wasmal.AudioBufferSourceNode
	bus          wasmal.AudioNode // due to WebAudio API limitations
	isPlaying    bool
}

var _ audio.PlaybackNode = (*PlaybackNode)(nil)
var _ Node = (*PlaybackNode)(nil)

func (n *PlaybackNode) AudioNode() wasmal.AudioNode {
	return n.bus
}

func (n *PlaybackNode) Start(offset float32) {
	if n.delegate != nil {
		n.Stop()
	}
	n.isPlaying = true
	n.delegate = n.audioContext.CreateBufferSource()
	n.delegate.SetBuffer(n.buffer)
	n.delegate.ConnectToNode(n.bus)
	n.delegate.StartOffset(0.0, float64(offset))
	n.delegate.SetOnEnded(n.onEnded)
}

func (n *PlaybackNode) Stop() {
	if n.delegate != nil {
		n.delegate.Stop(0)
		n.delegate.Disconnect()
		n.delegate = nil
	}
}

func (n *PlaybackNode) Resume() {

}

func (n *PlaybackNode) Pause() {

}

func (n *PlaybackNode) IsPlaying() bool {
	return n.isPlaying
}

func (n *PlaybackNode) IsLoop() bool {
	return n.delegate.Loop()
}

func (n *PlaybackNode) SetLoop(loop bool) {
	n.delegate.SetLoop(loop)
}

func (n *PlaybackNode) LoopStart() float32 {
	return float32(n.delegate.LoopStart())
}

func (n *PlaybackNode) SetLoopStart(loopStart float32) {
	n.delegate.SetLoopStart(float64(loopStart))
}

func (n *PlaybackNode) LoopEnd() float32 {
	return float32(n.delegate.LoopEnd())
}

func (n *PlaybackNode) SetLoopEnd(loopEnd float32) {
	n.delegate.SetLoopEnd(float64(loopEnd))
}

func (n *PlaybackNode) Delete() {
	if n.delegate != nil {
		n.delegate.Disconnect()
		n.delegate = nil
	}
}

func (n *PlaybackNode) onEnded() {
	n.isPlaying = false
	n.delegate.Disconnect()
	n.delegate = nil
}
