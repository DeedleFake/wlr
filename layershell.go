package wlr

// #include <wlr/types/wlr_layer_shell_v1.h>
import "C"

import "unsafe"

type LayerShellV1 struct {
	p *C.struct_wlr_layer_shell_v1
}

func NewLayerShellV1(display Display) LayerShellV1 {
	p := C.wlr_layer_shell_v1_create(display.p)
	return LayerShellV1{p: p}
}

func (ls LayerShellV1) OnDestroy(cb func(LayerShellV1)) Listener {
	return newListener(&ls.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(ls)
	})
}

func (ls LayerShellV1) OnNewSurface(cb func(LayerSurfaceV1)) Listener {
	return newListener(&ls.p.events.new_surface, func(lis Listener, data unsafe.Pointer) {
		cb(LayerSurfaceV1{p: (*C.struct_wlr_layer_surface_v1)(data)})
	})
}

type LayerSurfaceV1 struct {
	p *C.struct_wlr_layer_surface_v1
}
