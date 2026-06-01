package internal

import (
	"log/slog"

	"github.com/mokiat/lacking/audio"
	"github.com/mokiat/wasmal"
)

func NewPlayer() *Player {
	audioContext := wasmal.NewAudioContext()
	listener := newSpatialListener(audioContext)
	output := newOutputNode(audioContext)
	return &Player{
		audioContext: audioContext,
		listener:     listener,
		output:       output,
	}
}

type Player struct {
	audioContext wasmal.AudioContext
	listener     *SpatialListener
	output       *OutputNode
}

func (p *Player) SampleRate() int {
	return int(p.audioContext.SampleRate())
}

func (p *Player) CreateMedia(data audio.MediaData) *Media {
	frames := data.Frames
	if data.SampleRate != p.SampleRate() {
		logger.Warn("Resampling media",
			slog.Int("from", data.SampleRate),
			slog.Int("to", p.SampleRate()),
		)
		frames = audio.Resample(data.Frames, data.SampleRate, p.SampleRate())
	}

	buffer := p.audioContext.CreateBuffer(2, uint(len(frames)), uint(p.SampleRate()))
	if true {
		panic("TODO: Implement audio buffer data upload")
	}
	// TODO:
	// buffer.GetChannelData(0).Set(frames.Left)

	return &Media{
		buffer: buffer,
	}
}

func (p *Player) Output() *OutputNode {
	return p.output
}

func (p *Player) SpatialListener() *SpatialListener {
	return p.listener
}

func (p *Player) CreatePlaybackNode(media *Media) *PlaybackNode {
	return newPlaybackNode(p.audioContext, media)
}

func (p *Player) CreateOscillatorNode() *OscillatorNode {
	return newOscillatorNode(p.audioContext)
}

func (p *Player) CreateGainNode() *GainNode {
	return newGainNode(p.audioContext)
}

func (p *Player) CreatePanNode() *PanNode {
	return newPanNode(p.audioContext)
}

func (p *Player) CreateSpatialNode() *SpatialNode {
	return newSpatialNode(p.audioContext)
}

func (p *Player) CreateHighPassNode() *HighPassNode {
	return newHighPassNode(p.audioContext)
}

func (p *Player) CreateLowPassNode() *LowPassNode {
	return newLowPassNode(p.audioContext)
}

func (p *Player) CreateDelayNode() *DelayNode {
	return newDelayNode(p.audioContext)
}

func (p *Player) CreateReverbNode() *ReverbNode {
	return newReverbNode(p.audioContext)
}

func (p *Player) CreateCompressorNode() *CompressorNode {
	return newCompressorNode(p.audioContext)
}

func (p *Player) CreateConnectorNode() *ConnectorNode {
	return newConnectorNode(p.audioContext)
}

func (p *Player) Connect(from, to Node) {
	from.AudioNode().ConnectNode(to.AudioNode())
}

func (p *Player) Disconnect(from, to Node) {
	from.AudioNode().DisconnectNode(to.AudioNode())
}

func (p *Player) Play(media *Media, info audio.PlayInfo) *Playback {
	srcNode := p.CreatePlaybackNode(media)
	srcNode.SetLoop(info.Loop)
	srcNode.Start(0.0)

	panNode := p.CreatePanNode()
	panNode.SetPan(float32(info.Pan))

	gainNode := p.CreateGainNode()
	gainNode.SetGain(float32(info.Gain.ValueOrDefault(1.0)))

	p.Connect(srcNode, panNode)
	p.Connect(panNode, gainNode)
	p.Connect(gainNode, p.output)

	return &Playback{
		srcNode:  srcNode,
		panNode:  panNode,
		gainNode: gainNode,
	}
}

func (p *Player) Close() {
	p.audioContext.Close()
}
