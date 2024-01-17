package internal

import (
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/wasmgl"
)

func NewFence() *Fence {
	return &Fence{
		raw: wasmgl.FenceSync(wasmgl.SYNC_GPU_COMMANDS_COMPLETE, 0),
	}
}

type Fence struct {
	render.FenceMarker

	raw wasmgl.Sync
}

func (f *Fence) Status() render.FenceStatus {
	switch wasmgl.GetSyncParameter(f.raw, wasmgl.SYNC_STATUS).GLenum() {
	case wasmgl.SIGNALED:
		return render.FenceStatusSuccess
	case wasmgl.UNSIGNALED:
		return render.FenceStatusNotReady
	default:
		return render.FenceStatusNotReady
	}
}

func (f *Fence) Release() {
	wasmgl.DeleteSync(f.raw)
	f.raw = wasmgl.NilSync
}
