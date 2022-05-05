package wlr

// #include <wlr/types/wlr_layer_shell_v1.h>
import "C"

import "unsafe"

type LayerShellV1 struct {
	p *C.struct_wlr_layer_shell_v1
}

func NewLayerShellV1(display *Display) *LayerShellV1 {
	p := C.wlr_layer_shell_v1_create(display.p)
	trackObject(unsafe.Pointer(p), &p.events.destroy)
	return &LayerShellV1{p: p}
}

func (ls LayerShellV1) OnDestroy(cb func(LayerShellV1)) func() {
	lis := newListener(unsafe.Pointer(ls.p), func(lis *wlrlis, data unsafe.Pointer) {
		cb(ls)
	})
	C.wl_signal_add(&ls.p.events.destroy, lis)
	return func() {
		removeListener(lis)
	}
}

func (ls *LayerShellV1) OnNewSurface(cb func(*LayerSurfaceV1)) func() {
	lis := newListener(unsafe.Pointer(ls.p), func(lis *wlrlis, data unsafe.Pointer) {
		s := &LayerSurfaceV1{p: (*C.struct_wlr_layer_surface_v1)(data)}
		trackObject(data, &s.p.events.destroy)
		cb(s)
	})
	C.wl_signal_add(&ls.p.events.new_surface, lis)
	return func() {
		removeListener(lis)
	}
}

type LayerSurfaceV1 struct {
	p *C.struct_wlr_layer_surface_v1
}
