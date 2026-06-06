package internal

import (
	"math"

	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/lacking/core/audio"
	"github.com/mokiat/wasmal"
)

type PlaybackNode interface {
	Node
	InternalPause()
	InternalResume()
}

type BasePlayback struct {
	ctx wasmal.AudioContext

	sourceNode     wasmal.AudioBufferSourceNode
	gainFilter     *GainFilter
	lowPassFilter  *LowPassFilter
	highPassFilter *HighPassFilter

	callbackCleanup wasmal.CleanupFunc

	startContextTime float64

	startPosition float64
	pausePosition float64

	isInternalPaused bool
	isPaused         bool
	isPlaying        bool
}

var _ PlaybackNode = (*BasePlayback)(nil)

func NewBasePlayback(ctx wasmal.AudioContext, media *Media, settings audio.PlaybackSettings) *BasePlayback {
	sourceNode := ctx.CreateBufferSource()
	sourceNode.SetBuffer(media.buffer)
	lastOutput := wasmal.AudioNode(sourceNode)

	gainFilter := NewGainFilter(ctx)
	lastOutput.ConnectToNode(gainFilter.Input())
	lastOutput = gainFilter.Output()

	result := &BasePlayback{
		ctx:        ctx,
		sourceNode: sourceNode,
		gainFilter: gainFilter,
	}

	if settings.UseLowPassFilter {
		result.lowPassFilter = NewLowPassFilter(ctx)
		lastOutput.ConnectToNode(result.lowPassFilter.Input())
		lastOutput = result.lowPassFilter.Output()
	}

	if settings.UseHighPassFilter {
		result.highPassFilter = NewHighPassFilter(ctx)
		lastOutput.ConnectToNode(result.highPassFilter.Input())
		lastOutput = result.highPassFilter.Output()
	}

	return result
}

func (p *BasePlayback) Release() {
	if p.callbackCleanup != nil {
		p.callbackCleanup()
		p.callbackCleanup = nil
	}
}

func (p *BasePlayback) Input() wasmal.AudioNode {
	return p.sourceNode
}

func (p *BasePlayback) Output() wasmal.AudioNode {
	if p.highPassFilter != nil {
		return p.highPassFilter.Output()
	}
	if p.lowPassFilter != nil {
		return p.lowPassFilter.Output()
	}
	return p.gainFilter.Output()
}

func (p *BasePlayback) InternalPause() {
	if p.isInternalPaused {
		return
	}

	p.isInternalPaused = true

	if p.isPlaying && !p.isPaused {
		p.doPause()
	}
}

func (p *BasePlayback) InternalResume() {
	if !p.isInternalPaused {
		return
	}

	p.isInternalPaused = false

	if p.isPlaying && !p.isPaused {
		p.doResume()
	}
}

// Start begins playback from the given time offset in seconds.
func (p *BasePlayback) Start(at float64) {
	if p.isPlaying && !(p.isPaused || p.isInternalPaused) {
		p.doStop()
	}

	p.isPlaying = true
	p.isPaused = false

	if p.isInternalPaused {
		p.pausePosition = at // record desired position
		return
	}

	p.doStart(at)
}

// Stop halts playback and resets the position to the beginning.
func (p *BasePlayback) Stop() {
	if !p.isPlaying {
		return // not playing
	}
	p.isPlaying = false

	if p.isInternalPaused {
		return // source node is already inactive
	}

	p.doStop()
}

// Pause suspends playback without resetting the position.
func (p *BasePlayback) Pause() {
	if p.isPaused || !p.isPlaying {
		return // not playing or already paused
	}
	p.isPaused = true

	if p.isInternalPaused {
		return // source node is already inactive
	}

	p.doPause()
}

// Resume continues a previously paused playback.
func (p *BasePlayback) Resume() {
	if !p.isPaused {
		return // not paused
	}
	p.isPaused = false

	if !p.isPlaying {
		p.pausePosition = 0.0 // record desired position
		p.isPlaying = true    // force a start
	}

	if p.isInternalPaused {
		return // source node needs to remain inactive
	}

	p.doResume()
}

// Looping reports whether the playback loops when it reaches the end.
func (p *BasePlayback) Looping() bool {
	return p.sourceNode.Loop()
}

// SetLooping enables or disables looping.
func (p *BasePlayback) SetLooping(loop bool) {
	p.sourceNode.SetLoop(loop)
}

// LoopStart returns the loop start position in seconds.
func (p *BasePlayback) LoopStart() float64 {
	return p.sourceNode.LoopStart()
}

// SetLoopStart sets the loop start position in seconds.
func (p *BasePlayback) SetLoopStart(loopStart float64) {
	p.sourceNode.SetLoopStart(loopStart)
}

// LoopEnd returns the loop end position in seconds.
func (p *BasePlayback) LoopEnd() float64 {
	return p.sourceNode.LoopEnd()
}

// SetLoopEnd sets the loop end position in seconds.
func (p *BasePlayback) SetLoopEnd(loopEnd float64) {
	p.sourceNode.SetLoopEnd(loopEnd)
}

// Playing reports whether this playback is currently active.
func (p *BasePlayback) Playing() bool {
	return p.isPlaying
}

// PlaybackRate returns the playback speed multiplier; 1.0 is normal speed.
func (p *BasePlayback) PlaybackRate() float32 {
	return p.sourceNode.PlaybackRate().Value()
}

// SetPlaybackRate sets the playback speed multiplier.
func (p *BasePlayback) SetPlaybackRate(rate float32) {
	p.sourceNode.PlaybackRate().SetValue(rate)
}

// Gain returns the current gain of this playback.
func (p *BasePlayback) Gain() float32 {
	return p.gainFilter.Gain()
}

// SetGain sets the gain of this playback.
func (p *BasePlayback) SetGain(gain float32) {
	p.gainFilter.SetGain(gain)
}

// LowPassFilter returns the low-pass filter controls, or nil if the filter was
// not enabled at creation time.
func (p *BasePlayback) LowPassFilter() audio.FrequencyFilter {
	if p.lowPassFilter == nil {
		return nil
	}
	return p.lowPassFilter
}

// HighPassFilter returns the high-pass filter controls, or nil if the filter
// was not enabled at creation time.
func (p *BasePlayback) HighPassFilter() audio.FrequencyFilter {
	if p.highPassFilter == nil {
		return nil
	}
	return p.highPassFilter
}

// SetOnFinished sets a callback invoked when the playback reaches its end.
func (p *BasePlayback) SetOnFinished(onFinished func()) {
	if p.callbackCleanup != nil {
		p.callbackCleanup()
	}
	p.callbackCleanup = p.sourceNode.SetOnEnded(onFinished)
}

func (p *BasePlayback) doStop() {
	p.sourceNode.Stop(0.0)
}

func (p *BasePlayback) doStart(at float64) {
	p.startContextTime = p.ctx.CurrentTime()
	p.startPosition = at

	p.sourceNode.Disconnect()
	p.sourceNode = p.cloneSourceNode(p.sourceNode)
	p.sourceNode.ConnectToNode(p.gainFilter.Input())
	p.sourceNode.StartOffset(0.0, at)
}

func (p *BasePlayback) doPause() {
	// The logic here is simply horrid which huge likelyhood of bugs. All because
	// WebAudio API does not have a pause capability or a way to fetch the
	// current playback position.
	elapsedTime := p.ctx.CurrentTime() - p.startContextTime
	deltaPosition := elapsedTime * float64(p.PlaybackRate())
	absPosition := p.startPosition + deltaPosition
	if p.Looping() {
		relativePosition := absPosition - p.LoopStart()
		loopDuration := max(0.0001, p.LoopEnd()-p.LoopStart())
		p.pausePosition = p.LoopStart() + math.Mod(relativePosition, loopDuration)
	} else {
		p.pausePosition = absPosition
	}
	p.pausePosition = dprec.Clamp(p.pausePosition, 0.0, p.sourceNode.Buffer().Duration())

	p.sourceNode.Stop(0.0)
}

func (p *BasePlayback) doResume() {
	p.startContextTime = p.ctx.CurrentTime()
	p.startPosition = p.pausePosition

	p.sourceNode.Disconnect()
	p.sourceNode = p.cloneSourceNode(p.sourceNode)
	p.sourceNode.ConnectToNode(p.gainFilter.Input())
	p.sourceNode.StartOffset(0.0, p.pausePosition)
}

func (p *BasePlayback) cloneSourceNode(original wasmal.AudioBufferSourceNode) wasmal.AudioBufferSourceNode {
	result := p.ctx.CreateBufferSource()
	result.SetBuffer(original.Buffer())
	result.SetLoop(original.Loop())
	result.SetLoopStart(original.LoopStart())
	result.SetLoopEnd(original.LoopEnd())
	result.PlaybackRate().SetValue(original.PlaybackRate().Value())
	return result
}
