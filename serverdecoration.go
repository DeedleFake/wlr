package wlr

/*
#include <wlr/types/wlr_server_decoration.h>
*/
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

func CreateServerDecorationManager(display Display) ServerDecorationManager {
	p := C.wlr_server_decoration_manager_create(display.p)
	return ServerDecorationManager{p: p}
}

func (m ServerDecorationManager) OnDestroy(cb func(ServerDecorationManager)) Listener {
	return newListener(&m.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(m)
	})
}

func (m ServerDecorationManager) SetDefaultMode(mode ServerDecorationManagerMode) {
	C.wlr_server_decoration_manager_set_default_mode(m.p, C.uint32_t(mode))
}

func (m ServerDecorationManager) OnNewDecoration(cb func(ServerDecorationManager, ServerDecoration)) Listener {
	return newListener(&m.p.events.new_decoration, func(lis Listener, data unsafe.Pointer) {
		cb(
			m,
			ServerDecoration{p: (*C.struct_wlr_server_decoration)(data)},
		)
	})
}

func (d ServerDecoration) OnDestroy(cb func(ServerDecoration)) Listener {
	return newListener(&d.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(d)
	})
}

func (d ServerDecoration) OnMode(cb func(ServerDecoration)) Listener {
	return newListener(&d.p.events.mode, func(lis Listener, data unsafe.Pointer) {
		cb(d)
	})
}

func (d ServerDecoration) Mode() ServerDecorationManagerMode {
	return ServerDecorationManagerMode(d.p.mode)
}

func (d ServerDecoration) Surface() Surface {
	return Surface{p: d.p.surface}
}
