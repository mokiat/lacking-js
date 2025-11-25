package internal

import (
	"log/slog"

	"github.com/mokiat/lacking/audio"
	"github.com/mokiat/wasmal"
)

func NewPlayer() *Player {
	return &Player{
		audioContext: wasmal.NewAudioContext(),
	}
}

type Player struct {
	audioContext wasmal.AudioContext
}

func (p *Player) CreateMedia(info audio.MediaInfo) *Media {
	audioBufferPromise := p.audioContext.DecodeAudioData(info.Data)
	resolveAudioBuffer := make(chan wasmal.AudioBuffer, 1)
	resolveErr := make(chan error, 1)
	audioBufferPromise.Then(func(audioBuffer wasmal.AudioBuffer) {
		resolveAudioBuffer <- audioBuffer
	})
	audioBufferPromise.Catch(func(err error) {
		resolveErr <- err
	})
	select {
	case buffer := <-resolveAudioBuffer:
		return &Media{
			buffer: buffer,
		}
	case err := <-resolveErr:
		logger.Error("Error decoding media",
			slog.String("error", err.Error()),
		)
		return nil
	}
}

func (p *Player) Play(media *Media, info audio.PlayInfo) *Playback {
	gainNode := p.audioContext.CreateGain()
	gainNode.Gain().SetValue(info.Gain.ValueOrDefault(1.0))
	gainNode.ConnectNode(p.audioContext.Destination())

	panNode := p.audioContext.CreateStereoPanner()
	panNode.Pan().SetValue(info.Pan)
	panNode.ConnectNode(gainNode)

	bufferSource := p.audioContext.CreateBufferSource()
	bufferSource.SetBuffer(media.buffer)
	bufferSource.SetLoop(info.Loop)
	bufferSource.ConnectNode(panNode)
	bufferSource.Start(0.0)

	return &Playback{
		node: bufferSource,
	}
}

func (p *Player) Close() {
	p.audioContext.Close()
}
