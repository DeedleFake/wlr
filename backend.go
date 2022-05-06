package wlr

/*
#include <wayland-server-core.h>
#include <wlr/backend.h>
#include <wlr/backend/wayland.h>
#include <wlr/backend/x11.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

type Backend struct {
	p *C.struct_wlr_backend
}

func AutocreateBackend(display Display) Backend {
	p := C.wlr_backend_autocreate(display.p)
	return Backend{p: p}
}

func (b Backend) Destroy() {
	C.wlr_backend_destroy(b.p)
}

func (b Backend) OnDestroy(cb func(Backend)) Listener {
	return newListener(&b.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(b)
	})
}

func (b Backend) Start() error {
	if !C.wlr_backend_start(b.p) {
		return errors.New("can't start backend")
	}

	return nil
}

func (b Backend) OnNewOutput(cb func(Output)) Listener {
	return newListener(&b.p.events.new_output, func(lis Listener, data unsafe.Pointer) {
		cb(Output{p: (*C.struct_wlr_output)(data)})
	})
}

func (b Backend) OnNewInput(cb func(InputDevice)) Listener {
	return newListener(&b.p.events.new_input, func(lis Listener, data unsafe.Pointer) {
		cb(InputDevice{p: (*C.struct_wlr_input_device)(data)})
	})
}
