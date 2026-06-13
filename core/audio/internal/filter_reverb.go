package internal

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/core/audio"
	"github.com/mokiat/wasmal"
)

// ReverbFilter implements a Freeverb-style reverb using 8 parallel feedback
// comb filters (each with a DelayNode + damping BiquadFilter) followed by a
// ConvolverNode holding the precomputed impulse response of the series
// all-pass diffusion chain.
//
// Live feedback cycles are limited to the combs on purpose: Chrome processes
// audio graph cycles with an extra render quantum of loop latency, which
// merely detunes a comb slightly but breaks the exact feedforward/feedback
// cancellation a Schroeder all-pass depends on, producing audible resonances.
// The all-pass parameters are fixed, so their combined impulse response is
// generated once at construction and never needs to be rebuilt.
//
// All reverb parameters map to AudioParam values, so updates are cheap and
// require no graph rebuilds or buffer regeneration.
type ReverbFilter struct {
	input   wasmal.GainNode
	combs   [reverbCombFilterCount]*feedbackCombFilter
	dryGain wasmal.GainNode
	wetGain wasmal.GainNode
	output  wasmal.GainNode

	roomSize float32
	damping  float32
}

var _ audio.Reverb = (*ReverbFilter)(nil)
var _ Node = (*ReverbFilter)(nil)

func NewReverbFilter(ctx wasmal.AudioContext) *ReverbFilter {
	// Create all required nodes.
	input := ctx.CreateGain()

	var combs [reverbCombFilterCount]*feedbackCombFilter
	for i := range reverbCombFilterCount {
		comb := newFeedbackCombFilter(ctx)
		comb.SetFeedback(0.0)
		comb.SetDamping(0.0)
		comb.SetDelay(reverbCombFilterDelays[i])
		combs[i] = comb
	}

	combMix := ctx.CreateGain()
	combMix.Gain().SetValue(reverbCombFilterScale)

	diffusion := newDiffusionConvolver(ctx)

	dryGain := ctx.CreateGain()
	dryGain.Gain().SetValue(1.0)

	wetGain := ctx.CreateGain()
	wetGain.Gain().SetValue(0.5)

	output := ctx.CreateGain()

	// Wire all nodes.
	input.ConnectToNode(dryGain)

	for _, comb := range combs {
		input.ConnectToNode(comb.Input())
		comb.Output().ConnectToNode(combMix)
	}

	combMix.ConnectToNode(diffusion)
	diffusion.ConnectToNode(wetGain)

	dryGain.ConnectToNode(output)
	wetGain.ConnectToNode(output)

	result := &ReverbFilter{
		input:   input,
		combs:   combs,
		dryGain: dryGain,
		wetGain: wetGain,
		output:  output,
	}
	result.SetRoomSize(0.3)
	result.SetDamping(0.5)
	result.SetDry(1.0)
	result.SetWet(0.5)
	return result
}

func (f *ReverbFilter) Input() wasmal.AudioNode {
	return f.input
}

func (f *ReverbFilter) Output() wasmal.AudioNode {
	return f.output
}

func (f *ReverbFilter) RoomSize() float32 {
	return f.roomSize
}

func (f *ReverbFilter) SetRoomSize(size float32) {
	f.roomSize = sprec.Clamp(size, 0.0, 1.0)

	feedback := 0.3 + f.roomSize*0.6
	for _, comb := range f.combs {
		comb.SetFeedback(feedback)
	}
}

func (f *ReverbFilter) Damping() float32 {
	return f.damping
}

func (f *ReverbFilter) SetDamping(damping float32) {
	f.damping = sprec.Clamp(damping, 0.0, 1.0)

	for _, comb := range f.combs {
		comb.SetDamping(f.damping)
	}
}

func (f *ReverbFilter) Dry() float32 {
	return f.dryGain.Gain().Value()
}

func (f *ReverbFilter) SetDry(dry float32) {
	f.dryGain.Gain().SetValue(sprec.Clamp(dry, 0.0, 1.0))
}

func (f *ReverbFilter) Wet() float32 {
	return f.wetGain.Gain().Value()
}

func (f *ReverbFilter) SetWet(wet float32) {
	f.wetGain.Gain().SetValue(sprec.Clamp(wet, 0.0, 1.0))
}

const (
	reverbCombFilterCount = 8
	reverbCombFilterScale = 1.0 / float32(reverbCombFilterCount)

	reverbDiffusionIRSeconds    = 0.25
	reverbAllPassFilterCount    = 4
	reverbAllPassFilterFeedback = 0.5
	reverbAllPassFilterSpread   = 0.000522 // ~23 samples @ 44100 Hz
)

var (
	// Delay times back-calculated from Freeverb's original sample counts at 44100 Hz.
	// They are stored as seconds so they scale correctly at other sample rates.
	reverbCombFilterDelays = [reverbCombFilterCount]float32{
		0.025306, // 1116 samples @ 44100 Hz
		0.026939, // 1188 samples @ 44100 Hz
		0.028957, // 1277 samples @ 44100 Hz
		0.030748, // 1356 samples @ 44100 Hz
		0.032245, // 1422 samples @ 44100 Hz
		0.033810, // 1491 samples @ 44100 Hz
		0.035306, // 1557 samples @ 44100 Hz
		0.036667, // 1617 samples @ 44100 Hz
	}
	reverbAllPassFilterDelays = [reverbAllPassFilterCount]float64{
		0.012608, // 556 samples @ 44100 Hz
		0.010000, // 441 samples @ 44100 Hz
		0.007732, // 341 samples @ 44100 Hz
		0.005102, // 225 samples @ 44100 Hz
	}
)

type feedbackCombFilter struct {
	sampleRate float32
	input      wasmal.GainNode
	delay      wasmal.DelayNode
	damping    wasmal.BiquadFilterNode
	feedback   wasmal.GainNode
}

func newFeedbackCombFilter(ctx wasmal.AudioContext) *feedbackCombFilter {
	input := ctx.CreateGain()

	delay := ctx.CreateDelay(1.0)

	damping := ctx.CreateBiquadFilter()
	damping.SetType(wasmal.BiquadFilterTypeLowpass)
	damping.Q().SetValue(0.5)

	feedback := ctx.CreateGain()

	input.ConnectToNode(delay)
	delay.ConnectToNode(damping)
	damping.ConnectToNode(feedback)
	feedback.ConnectToNode(delay)

	return &feedbackCombFilter{
		sampleRate: ctx.SampleRate(),

		input:    input,
		delay:    delay,
		damping:  damping,
		feedback: feedback,
	}
}

func (f *feedbackCombFilter) Input() wasmal.AudioNode {
	return f.input
}

func (f *feedbackCombFilter) Output() wasmal.AudioNode {
	return f.delay
}

func (f *feedbackCombFilter) SetFeedback(feedback float32) {
	f.feedback.Gain().SetValue(feedback)
}

func (f *feedbackCombFilter) SetDamping(damping float32) {
	frequency := ((1.0 - damping) * f.sampleRate) / (sprec.Tau * damping)
	frequency = sprec.Clamp(frequency, 0.1, f.sampleRate/2.0)
	f.damping.Frequency().SetValue(frequency)
}

func (f *feedbackCombFilter) SetDelay(delaySeconds float32) {
	f.delay.DelayTime().SetValue(delaySeconds)
}

// newDiffusionConvolver creates a ConvolverNode loaded with the impulse
// response of the four series all-pass filters, with the right channel using
// spread delays for stereo decorrelation.
func newDiffusionConvolver(ctx wasmal.AudioContext) wasmal.ConvolverNode {
	sampleRate := int(ctx.SampleRate())
	left := diffusionImpulseResponse(sampleRate, 0.0)
	right := diffusionImpulseResponse(sampleRate, reverbAllPassFilterSpread)

	buffer := ctx.CreateBuffer(2, uint32(len(left)), ctx.SampleRate())
	buffer.GetChannelData(0).CopyFrom(left)
	buffer.GetChannelData(1).CopyFrom(right)

	convolver := ctx.CreateConvolver()
	// Normalization must be disabled (before assigning the buffer) so that
	// the diffusion stage preserves the unity gain of the all-pass chain.
	convolver.SetNormalize(false)
	convolver.SetBuffer(buffer)
	return convolver
}

func diffusionImpulseResponse(sampleRate int, extraDelay float64) []float32 {
	var filters [reverbAllPassFilterCount]*allPassSampleFilter
	for i := range reverbAllPassFilterCount {
		delay := audio.SampleCount(reverbAllPassFilterDelays[i]+extraDelay, sampleRate)
		filters[i] = newAllPassSampleFilter(delay)
	}

	result := make([]float32, audio.SampleCount(reverbDiffusionIRSeconds, sampleRate))
	for n := range result {
		sample := float32(0.0)
		if n == 0 {
			sample = 1.0
		}
		for _, filter := range filters {
			sample = filter.ProcessSample(sample)
		}
		result[n] = sample
	}
	return result
}

// allPassSampleFilter implements a Schroeder all-pass filter on individual
// samples. It mirrors the native implementation and is used only to generate
// the diffusion impulse response.
type allPassSampleFilter struct {
	buffer []float32
	delay  int
	write  int
}

func newAllPassSampleFilter(delaySamples int) *allPassSampleFilter {
	return &allPassSampleFilter{
		buffer: make([]float32, delaySamples+1),
		delay:  delaySamples,
	}
}

func (f *allPassSampleFilter) ProcessSample(input float32) float32 {
	size := len(f.buffer)
	read := (f.write - f.delay + size) % size
	buffered := f.buffer[read]
	output := buffered - input*reverbAllPassFilterFeedback
	f.buffer[f.write] = input + output*reverbAllPassFilterFeedback
	f.write = (f.write + 1) % size
	return output
}
