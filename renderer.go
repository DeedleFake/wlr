package wlr

/*
#include <wlr/render/wlr_renderer.h>
*/
import "C"

import (
	"image"
	"image/color"
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

func (r Renderer) Begin(output Output, width int, height int) {
	C.wlr_renderer_begin(r.p, C.uint(width), C.uint(height))
}

func (r Renderer) Clear(c color.Color) {
	cc := colorToC(c)
	C.wlr_renderer_clear(r.p, &cc[0])
}

func (r Renderer) End() {
	C.wlr_renderer_end(r.p)
}

func (r Renderer) RenderTextureWithMatrix(texture Texture, matrix *Matrix, alpha float32) {
	m := matrix.toC()
	C.wlr_render_texture_with_matrix(r.p, texture.p, &m[0], C.float(alpha))
}

func (r *Renderer) RenderRect(box image.Rectangle, c color.Color, projection *Matrix) {
	cc := colorToC(c)
	pm := projection.toC()
	C.wlr_render_rect(
		r.p,
		boxToC(box),
		&cc[0],
		&pm[0],
	)
}
