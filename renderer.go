package wlr

/*
#include <wlr/render/wlr_renderer.h>
*/
import "C"

import (
	"unsafe"
)

type Renderer struct {
	p *C.struct_wlr_renderer
}

func AutocreateRenderer(backend Backend) Renderer {
	p := C.wlr_renderer_autocreate(backend.p)
	return Renderer{p: p}
}

func (r Renderer) Destroy() {
	C.wlr_renderer_destroy(r.p)
}

func (r Renderer) OnDestroy(cb func(Renderer)) Listener {
	return newListener(&r.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(r)
	})
}

func (r Renderer) InitWLDisplay(display Display) {
	C.wlr_renderer_init_wl_display(r.p, display.p)
}

func (r Renderer) InitWLSHM(display Display) {
	C.wlr_renderer_init_wl_shm(r.p, display.p)
}
