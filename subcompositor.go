package wlr

/*
#include <wlr/types/wlr_subcompositor.h>
*/
import "C"

type Subcompositor struct {
	p *C.struct_wlr_subcompositor
}

func SubcompositorCreate(display Display) Subcompositor {
	p := C.wlr_subcompositor_create(display.p)
	return Subcompositor{p: p}
}
