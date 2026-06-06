package audio

import (
	"github.com/mokiat/lacking-js/core/audio/internal"
	"github.com/mokiat/lacking/core/audio"
	"github.com/mokiat/wasmal"
)

type API struct {
	ctx wasmal.AudioContext

	masterBus *internal.MasterBus
	listener  *internal.SpatialListener
}

var _ audio.API = (*API)(nil)

func NewAPI() *API {
	ctx := wasmal.NewAudioContext()

	masterBus := internal.NewMasterBus(ctx)
	masterBus.Output().ConnectToNode(ctx.Destination())

	return &API{
		ctx:       ctx,
		masterBus: masterBus,
		listener:  internal.NewSpatialListener(ctx),
	}
}

func (a *API) CreateMedia(data audio.MediaData) audio.Media {
	return internal.NewMedia(a.ctx, data)
}

func (a *API) CreateBus(settings audio.BusSettings) audio.Bus {
	bus := internal.NewBus(a.ctx, settings)
	bus.Output().ConnectToNode(a.masterBus.Input())
	return bus
}

func (a *API) CreatePlayback(targetBus audio.Bus, targetMedia audio.Media, settings audio.PlaybackSettings) audio.Playback {
	bus := targetBus.(*internal.Bus)
	media := targetMedia.(*internal.Media)

	basePlayback := internal.NewBasePlayback(
		a.ctx,
		media,
		settings,
	)
	playback := internal.NewDefaultPlayback(basePlayback)

	bus.AddPlayback(playback)

	return playback
}

func (a *API) CreateSpatialPlayback(targetBus audio.Bus, targetMedia audio.Media, settings audio.PlaybackSettings) audio.SpatialPlayback {
	bus := targetBus.(*internal.Bus)
	media := targetMedia.(*internal.Media)

	basePlayback := internal.NewBasePlayback(
		a.ctx,
		media,
		settings,
	)
	playback := internal.NewSpatialPlayback(basePlayback)

	bus.AddPlayback(playback)

	return playback
}

func (a *API) MasterBus() audio.MasterBus {
	return a.masterBus
}

func (a *API) SpatialListener() audio.SpatialListener {
	return a.listener
}
