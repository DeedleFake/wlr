package wlr

/*
#include <wlr/types/wlr_xdg_decoration_v1.h>
*/
import "C"

import "unsafe"

type XDGToplevelDecorationV1Mode uint32

const (
	XDGToplevelDecorationV1ModeNone       XDGToplevelDecorationV1Mode = C.WLR_XDG_TOPLEVEL_DECORATION_V1_MODE_NONE
	XDGToplevelDecorationV1ModeClientSide XDGToplevelDecorationV1Mode = C.WLR_XDG_TOPLEVEL_DECORATION_V1_MODE_CLIENT_SIDE
	XDGToplevelDecorationV1ModeServerSide XDGToplevelDecorationV1Mode = C.WLR_XDG_TOPLEVEL_DECORATION_V1_MODE_SERVER_SIDE
)

type XDGDecorationManagerV1 struct {
	p *C.struct_wlr_xdg_decoration_manager_v1
}

func CreateXDGDecorationManagerV1(display Display) XDGDecorationManagerV1 {
	p := C.wlr_xdg_decoration_manager_v1_create(display.p)
	return XDGDecorationManagerV1{p: p}
}

func (m XDGDecorationManagerV1) OnNewToplevelDecoration(cb func(d XDGToplevelDecorationV1)) Listener {
	return newListener(&m.p.events.new_toplevel_decoration, func(lis Listener, data unsafe.Pointer) {
		cb(XDGToplevelDecorationV1{p: (*C.struct_wlr_xdg_toplevel_decoration_v1)(data)})
	})
}

func (m XDGDecorationManagerV1) OnDestroy(cb func()) Listener {
	return newListener(&m.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb()
	})
}

type XDGToplevelDecorationV1 struct {
	p *C.struct_wlr_xdg_toplevel_decoration_v1
}

func (d XDGToplevelDecorationV1) RequestedMode() XDGToplevelDecorationV1Mode {
	return XDGToplevelDecorationV1Mode(d.p.requested_mode)
}

func (d XDGToplevelDecorationV1) SetMode(mode XDGToplevelDecorationV1Mode) {
	C.wlr_xdg_toplevel_decoration_v1_set_mode(d.p, C.enum_wlr_xdg_toplevel_decoration_v1_mode(mode))
}

func (d XDGToplevelDecorationV1) OnDestroy(cb func()) Listener {
	return newListener(&d.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb()
	})
}

type XDGToplevelDecorationV1State struct {
	v C.struct_wlr_xdg_toplevel_decoration_v1_state
}

func (s XDGToplevelDecorationV1State) Mode() XDGToplevelDecorationV1Mode {
	return XDGToplevelDecorationV1Mode(s.v.mode)
}
