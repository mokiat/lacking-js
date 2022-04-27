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
	render.FenceObject
	raw wasmgl.Sync
}

func (f *Fence) Status() render.FenceStatus {
	switch wasmgl.ClientWaitSync(f.raw, wasmgl.SYNC_FLUSH_COMMANDS_BIT, 0) {
	case wasmgl.ALREADY_SIGNALED, wasmgl.CONDITION_SATISFIED:
		return render.FenceStatusSuccess
	case wasmgl.TIMEOUT_EXPIRED:
		return render.FenceStatusNotReady
	default:
		return render.FenceStatusDeviceLost
	}
}

func (f *Fence) Delete() {
	wasmgl.DeleteSync(f.raw)
	f.raw = wasmgl.NilSync
}
