package wlr

// #include <wayland-server-core.h>
// #include <wlr/backend.h>
// #include <wlr/backend/wayland.h>
// #include <wlr/backend/x11.h>
import "C"

import (
	"errors"
	"unsafe"
)

type Backend struct {
	p *C.struct_wlr_backend
}

func NewBackend(display *Display) *Backend {
	p := C.wlr_backend_autocreate(display.p)
	trackObject(unsafe.Pointer(p), &p.events.destroy)
	return &Backend{p: p}
}

func (b *Backend) Destroy() {
	C.wlr_backend_destroy(b.p)
}

func (b *Backend) OnDestroy(cb func(*Backend)) func() {
	lis := newListener(unsafe.Pointer(b.p), func(*wlrlis, unsafe.Pointer) {
		cb(b)
	})
	C.wl_signal_add(&b.p.events.destroy, lis)
	return func() {
		removeListener(lis)
	}
}

func (b *Backend) Start() error {
	if !C.wlr_backend_start(b.p) {
		return errors.New("can't start backend")
	}

	return nil
}

func (b *Backend) OnNewOutput(cb func(*Output)) func() {
	lis := newListener(unsafe.Pointer(b.p), func(lis *wlrlis, data unsafe.Pointer) {
		output := &Output{p: (*C.struct_wlr_output)(data)}
		trackObject(data, &output.p.events.destroy)
		cb(output)
	})
	C.wl_signal_add(&b.p.events.new_output, lis)
	return func() {
		removeListener(lis)
	}
}

func (b *Backend) OnNewInput(cb func(*InputDevice)) func() {
	lis := newListener(unsafe.Pointer(b.p), func(lis *wlrlis, data unsafe.Pointer) {
		dev := &InputDevice{p: (*C.struct_wlr_input_device)(data)}
		trackObject(data, &dev.p.events.destroy)
		cb(dev)
	})
	C.wl_signal_add(&b.p.events.new_input, lis)
	return func() {
		removeListener(lis)
	}
}
