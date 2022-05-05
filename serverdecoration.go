package wlr

// #include <wlr/types/wlr_server_decoration.h>
import "C"

import "unsafe"

type ServerDecorationManagerMode uint32

const (
	ServerDecorationManagerModeNone   ServerDecorationManagerMode = C.WLR_SERVER_DECORATION_MANAGER_MODE_NONE
	ServerDecorationManagerModeClient ServerDecorationManagerMode = C.WLR_SERVER_DECORATION_MANAGER_MODE_CLIENT
	ServerDecorationManagerModeServer ServerDecorationManagerMode = C.WLR_SERVER_DECORATION_MANAGER_MODE_SERVER
)

type ServerDecorationManager struct {
	p *C.struct_wlr_server_decoration_manager
}

type ServerDecoration struct {
	p *C.struct_wlr_server_decoration
}

func NewServerDecorationManager(display *Display) *ServerDecorationManager {
	p := C.wlr_server_decoration_manager_create(display.p)
	trackObject(unsafe.Pointer(p), &p.events.destroy)
	return &ServerDecorationManager{p: p}
}

func (m *ServerDecorationManager) OnDestroy(cb func(*ServerDecorationManager)) func() {
	lis := newListener(unsafe.Pointer(m.p), func(lis *wlrlis, data unsafe.Pointer) {
		cb(m)
	})
	C.wl_signal_add(&m.p.events.destroy, lis)
	return func() {
		removeListener(lis)
	}
}

func (m ServerDecorationManager) SetDefaultMode(mode ServerDecorationManagerMode) {
	C.wlr_server_decoration_manager_set_default_mode(m.p, C.uint32_t(mode))
}

func (m *ServerDecorationManager) OnNewMode(cb func(*ServerDecorationManager, *ServerDecoration)) func() {
	lis := newListener(unsafe.Pointer(m.p), func(lis *wlrlis, data unsafe.Pointer) {
		dec := &ServerDecoration{
			p: (*C.struct_wlr_server_decoration)(data),
		}
		trackObject(unsafe.Pointer(dec.p), &dec.p.events.destroy)
		cb(m, dec)
	})
	C.wl_signal_add(&m.p.events.new_decoration, lis)
	return func() {
		removeListener(lis)
	}
}

func (d *ServerDecoration) OnDestroy(cb func(*ServerDecoration)) func() {
	lis := newListener(unsafe.Pointer(d.p), func(lis *wlrlis, data unsafe.Pointer) {
		cb(d)
	})
	C.wl_signal_add(&d.p.events.destroy, lis)
	return func() {
		removeListener(lis)
	}
}

func (d *ServerDecoration) OnMode(cb func(*ServerDecoration)) func() {
	lis := newListener(unsafe.Pointer(d.p), func(lis *wlrlis, data unsafe.Pointer) {
		cb(d)
	})
	C.wl_signal_add(&d.p.events.mode, lis)
	return func() {
		removeListener(lis)
	}
}

func (d ServerDecoration) Mode() ServerDecorationManagerMode {
	return ServerDecorationManagerMode(d.p.mode)
}
