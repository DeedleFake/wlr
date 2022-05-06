package wlr

/*
#include <wlr/types/wlr_screencopy_v1.h>
*/
import "C"

import "unsafe"

type ScreencopyManagerV1 struct {
	p *C.struct_wlr_screencopy_manager_v1
}

func CreateScreencopyManagerV1(display Display) ScreencopyManagerV1 {
	p := C.wlr_screencopy_manager_v1_create(display.p)
	return ScreencopyManagerV1{p: p}
}

func (m ScreencopyManagerV1) OnDestroy(cb func(ScreencopyManagerV1)) Listener {
	return newListener(&m.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(m)
	})
}
