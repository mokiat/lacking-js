package internal

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/core/audio"
	"github.com/mokiat/wasmal"
)

// ReverbFilter implements a Freeverb-style reverb using 8 parallel feedback
// comb filters (each with a DelayNode + damping BiquadFilter) followed by 4
// series all-pass BiquadFilterNodes for diffusion. Dry/wet balance is
// controlled via GainNodes without requiring any graph rebuilds.
type ReverbFilter struct {
	ctx wasmal.AudioContext

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
		comb.SetDelay(reverbCombFilterDelays[i])
		combs[i] = comb
	}

	combMix := ctx.CreateGain()
	combMix.Gain().SetValue(reverbCombFilterScale)

	var allPasses [reverbAllPassFilterCount]*allPassFilter
	for i := range reverbAllPassFilterCount {
		allPass := newAllPassFilter(ctx)
		allPass.SetFeedback(reverbAllPassFilterFeedback)
		allPass.SetDelay(reverbAllPassFilterDelays[i])
		allPasses[i] = allPass
	}

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

	combMix.ConnectToNode(allPasses[0].Input())
	for i := 1; i < reverbAllPassFilterCount; i++ {
		prevAllPass := allPasses[i-1]
		currentAllPass := allPasses[i]
		prevAllPass.Output().ConnectToNode(currentAllPass.Input())
	}
	allPasses[reverbAllPassFilterCount-1].Output().ConnectToNode(wetGain)

	dryGain.ConnectToNode(output)
	wetGain.ConnectToNode(output)

	result := &ReverbFilter{
		ctx:     ctx,
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
	reverbCombFilterCount       = 8
	reverbCombFilterScale       = 1.0 / float32(reverbCombFilterCount)
	reverbAllPassFilterCount    = 4
	reverbAllPassFilterFeedback = 0.5
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
	reverbAllPassFilterDelays = [reverbAllPassFilterCount]float32{
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

type allPassFilter struct {
	input    wasmal.GainNode
	direct   wasmal.GainNode
	delay    wasmal.DelayNode
	feedback wasmal.GainNode
	output   wasmal.GainNode
}

func newAllPassFilter(ctx wasmal.AudioContext) *allPassFilter {
	input := ctx.CreateGain()

	direct := ctx.CreateGain()
	direct.Gain().SetValue(-0.5)

	delay := ctx.CreateDelay(1.0)
	delay.DelayTime().SetValue(0.01)

	feedback := ctx.CreateGain()
	feedback.Gain().SetValue(0.5)

	output := ctx.CreateGain()

	input.ConnectToNode(direct)
	direct.ConnectToNode(output)

	input.ConnectToNode(delay)
	delay.ConnectToNode(output)
	output.ConnectToNode(feedback)
	feedback.ConnectToNode(delay)

	return &allPassFilter{
		input:    input,
		direct:   direct,
		delay:    delay,
		feedback: feedback,
		output:   output,
	}
}

func (f *allPassFilter) Input() wasmal.AudioNode {
	return f.input
}

func (f *allPassFilter) Output() wasmal.AudioNode {
	return f.output
}

func (f *allPassFilter) SetFeedback(feedback float32) {
	f.direct.Gain().SetValue(-feedback)
	f.feedback.Gain().SetValue(feedback)
}

func (f *allPassFilter) SetDelay(delaySeconds float32) {
	f.delay.DelayTime().SetValue(delaySeconds)
}
