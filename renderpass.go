package wlr

/*
#include <wlr/render/pass.h>

// options can't hold a Go pointer when it's passed to C, so alpha is
// pointed to from here instead.
static inline void _wlr_render_pass_add_texture(struct wlr_render_pass *pass, struct wlr_render_texture_options *options, float alpha) {
	options->alpha = &alpha;
	wlr_render_pass_add_texture(pass, options);
}
*/
import "C"

import (
	"image"
	"image/color"
)

type BlendMode uint32

const (
	BlendModePremultiplied BlendMode = C.WLR_RENDER_BLEND_MODE_PREMULTIPLIED
	BlendModeNone          BlendMode = C.WLR_RENDER_BLEND_MODE_NONE
)

type FilterMode uint32

const (
	FilterBilinear FilterMode = C.WLR_SCALE_FILTER_BILINEAR
	FilterNearest  FilterMode = C.WLR_SCALE_FILTER_NEAREST
)

type RenderPass struct {
	p *C.struct_wlr_render_pass
}

func (r RenderPass) Valid() bool {
	return r.p != nil
}

func (r RenderPass) Submit() bool {
	return bool(C.wlr_render_pass_submit(r.p))
}

func (r RenderPass) AddTexture(texture Texture, srcBox image.Rectangle, dstBox image.Rectangle, alpha float32, transform OutputTransform, filterMode FilterMode, blendMode BlendMode) {
	var options C.struct_wlr_render_texture_options
	options.texture = texture.p
	if srcBox != (image.Rectangle{}) {
		options.src_box = fboxFromRect(srcBox)
	}
	if dst := boxToC(dstBox); dst != nil {
		options.dst_box = *dst
	}
	options.transform = C.enum_wl_output_transform(transform)
	options.filter_mode = C.enum_wlr_scale_filter_mode(filterMode)
	options.blend_mode = C.enum_wlr_render_blend_mode(blendMode)
	C._wlr_render_pass_add_texture(r.p, &options, C.float(alpha))
}

func (r RenderPass) AddRect(box image.Rectangle, c color.Color, blendMode BlendMode) {
	cc := colorToC(c)
	var options C.struct_wlr_render_rect_options
	if b := boxToC(box); b != nil {
		options.box = *b
	}
	options.color = C.struct_wlr_render_color{r: cc[0], g: cc[1], b: cc[2], a: cc[3]}
	options.blend_mode = C.enum_wlr_render_blend_mode(blendMode)
	C.wlr_render_pass_add_rect(r.p, &options)
}

func fboxFromRect(r image.Rectangle) C.struct_wlr_fbox {
	r = r.Canon()
	return C.struct_wlr_fbox{
		x:      C.double(r.Min.X),
		y:      C.double(r.Min.Y),
		width:  C.double(r.Dx()),
		height: C.double(r.Dy()),
	}
}
