package wlr

/*
#include <wlr/types/wlr_presentation_time.h>
*/
import "C"

import "unsafe"

type Presentation struct {
	p *C.struct_wlr_presentation
}

// CreatePresentation creates the wp_presentation global. Feedback for a
// surface is only sent once the compositor reports the output that the
// surface was shown on with PresentationSurfaceTexturedOnOutput or
// PresentationSurfaceScannedOutOnOutput.
func CreatePresentation(display Display, backend Backend, version uint32) Presentation {
	p := C.wlr_presentation_create(display.p, backend.p, C.uint32_t(version))
	return Presentation{p: p}
}

func (p Presentation) OnDestroy(cb func(Presentation)) Listener {
	return newListener(&p.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(p)
	})
}

// PresentationSurfaceTexturedOnOutput tells clients waiting on
// presentation feedback for surface that it was rendered on output as a
// texture. Call it before committing the output.
func PresentationSurfaceTexturedOnOutput(surface Surface, output Output) {
	C.wlr_presentation_surface_textured_on_output(surface.p, output.p)
}

// PresentationSurfaceScannedOutOnOutput is like
// PresentationSurfaceTexturedOnOutput, but for a surface whose buffer
// was scanned out directly.
func PresentationSurfaceScannedOutOnOutput(surface Surface, output Output) {
	C.wlr_presentation_surface_scanned_out_on_output(surface.p, output.p)
}
