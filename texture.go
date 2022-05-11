package wlr

/*
#include <wlr/render/wlr_texture.h>
*/
import "C"

import (
	"image"
	"unsafe"

	"golang.org/x/image/draw"
)

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

func TextureFromImage(renderer Renderer, img image.Image) Texture {
	nrgba, ok := img.(*image.NRGBA)
	if !ok {
		nrgba = image.NewNRGBA(img.Bounds())
		draw.Copy(nrgba, image.ZP, img, nrgba.Bounds(), draw.Src, nil)
	}

	return TextureFromPixels(
		renderer,
		uint32('A'|('R'<<8)|('2'<<16)|('4'<<24)),
		uint32(nrgba.Stride),
		uint32(nrgba.Bounds().Dx()),
		uint32(nrgba.Bounds().Dy()),
		nrgba.Pix,
	)
}

func (t Texture) Destroy() {
	C.wlr_texture_destroy(t.p)
}

func (t Texture) Valid() bool {
	return t.p != nil
}

func (t Texture) Width() int {
	return int(t.p.width)
}

func (t Texture) Height() int {
	return int(t.p.height)
}
