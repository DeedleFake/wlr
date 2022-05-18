package wlr

/*
#include <wlr/types/wlr_buffer.h>
#include <wlr/types/wlr_drm.h>
*/
import "C"

type DRM struct {
	p *C.struct_wlr_drm
}

func CreateDRM(d Display, r Renderer) DRM {
	p := C.wlr_drm_create(d.p, r.p)
	return DRM{p: p}
}
