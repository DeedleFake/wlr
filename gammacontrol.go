package wlr

import "unsafe"

type GammaControlManagerV1 struct {
	p *C.struct_wlr_gamma_control_manager_v1
}

func CreateGammaControlManagerV1(display Display) GammaControlManagerV1 {
	p := C.wlr_gamma_control_manager_v1_create(display.p)
	return GammaControlManagerV1{p: p}
}

func (m GammaControlManagerV1) OnDestroy(cb func(GammaControlManagerV1)) Listener {
	return newListener(&m.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(m)
	})
}
