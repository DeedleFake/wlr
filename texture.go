package wlr

/*
#include <wlr/render/wlr_texture.h>
*/
import "C"

import "unsafe"

type Texture struct {
	p *C.struct_wlr_texture
}

func TextureFromPixels(renderer Renderer, fmt, stride, width, height uint32, data []byte) Texture {
	p := C.wlr_texture_from_pixels(
		renderer.p,
		C.uint32_t(fmt),
		C.uint32_t(stride),
		C.uint32_t(width),
		C.uint32_t(height),
		unsafe.Pointer(&data[0]), // TODO: Does this need to be allocated by C?
	)
	return Texture{p: p}
}

func (t Texture) Destroy() {
	C.wlr_texture_destroy(t.p)
}

func (t Texture) Valid() bool {
	return t.p == nil
}
