package wlr

// #include <wlr/backend.h>
// #include <wlr/types/wlr_xdg_shell.h>
//
// void _wlr_xdg_surface_for_each_cb(struct wlr_surface *surface, int sx, int sy, void *data);
//
// static inline void _wlr_xdg_surface_for_each_surface(struct wlr_xdg_surface *surface, void *user_data) {
// 	wlr_xdg_surface_for_each_surface(surface, &_wlr_xdg_surface_for_each_cb, user_data);
// }
import "C"

import (
	"sync"
	"unsafe"
)

type XDGSurfaceRole uint32

const (
	XDGSurfaceRoleNone     XDGSurfaceRole = C.WLR_XDG_SURFACE_ROLE_NONE
	XDGSurfaceRoleTopLevel XDGSurfaceRole = C.WLR_XDG_SURFACE_ROLE_TOPLEVEL
	XDGSurfaceRolePopup    XDGSurfaceRole = C.WLR_XDG_SURFACE_ROLE_POPUP
)

var (
	xdgSurfaceWalkers      = map[*C.struct_wlr_xdg_surface]XDGSurfaceWalkFunc{}
	xdgSurfaceWalkersMutex sync.RWMutex
)

type XDGSurfaceWalkFunc func(surface Surface, sx int, sy int)

type XDGShell struct {
	p *C.struct_wlr_xdg_shell
}

func NewXDGShell(display *Display) *XDGShell {
	p := C.wlr_xdg_shell_create(display.p)
	trackObject(unsafe.Pointer(p), &p.events.destroy)
	return &XDGShell{p: p}
}

func (s *XDGShell) OnDestroy(cb func(*XDGShell)) func() {
	lis := newListener(unsafe.Pointer(s.p), func(lis *wlrlis, data unsafe.Pointer) {
		cb(s)
	})
	C.wl_signal_add(&s.p.events.destroy, lis)
	return func() {
		removeListener(lis)
	}
}

func (s *XDGShell) OnNewSurface(cb func(*XDGSurface)) func() {
	lis := newListener(unsafe.Pointer(s.p), func(lis *wlrlis, data unsafe.Pointer) {
		surface := &XDGSurface{p: (*C.struct_wlr_xdg_surface)(data)}
		trackObject(data, &surface.p.events.destroy)
		//man.add(unsafe.Pointer(surface.p), &surface.p.events.destroy, func(data unsafe.Pointer) {
		//	man.delete(unsafe.Pointer(surface.p))
		//	man.delete(unsafe.Pointer(surface.TopLevel().p))
		//})
		//man.add(unsafe.Pointer(surface.p.surface), &surface.p.surface.events.destroy, func(data unsafe.Pointer) {
		//	man.delete(unsafe.Pointer(surface.p.surface))
		//})
		cb(surface)
	})
	C.wl_signal_add(&s.p.events.new_surface, lis)
	return func() {
		removeListener(lis)
	}
}

type XDGSurface struct {
	p *C.struct_wlr_xdg_surface
}

func (s XDGSurface) Nil() bool {
	return s.p == nil
}

func (s XDGSurface) Walk(visit XDGSurfaceWalkFunc) {
	xdgSurfaceWalkersMutex.Lock()
	xdgSurfaceWalkers[s.p] = visit
	xdgSurfaceWalkersMutex.Unlock()

	C._wlr_xdg_surface_for_each_surface(s.p, unsafe.Pointer(s.p))

	xdgSurfaceWalkersMutex.Lock()
	delete(xdgSurfaceWalkers, s.p)
	xdgSurfaceWalkersMutex.Unlock()
}

func (s XDGSurface) Role() XDGSurfaceRole {
	return XDGSurfaceRole(s.p.role)
}

func (s XDGSurface) TopLevel() XDGTopLevel {
	p := *(*unsafe.Pointer)(unsafe.Pointer(&s.p.anon0[0]))
	return XDGTopLevel{p: (*C.struct_wlr_xdg_toplevel)(p)}
}

func (s XDGSurface) TopLevelSetActivated(activated bool) {
	C.wlr_xdg_toplevel_set_activated(s.p, C.bool(activated))
}

func (s XDGSurface) TopLevelSetSize(width uint32, height uint32) {
	C.wlr_xdg_toplevel_set_size(s.p, C.uint32_t(width), C.uint32_t(height))
}

func (s XDGSurface) TopLevelSetTiled(edges Edges) {
	C.wlr_xdg_toplevel_set_tiled(s.p, C.uint32_t(edges))
}

func (s XDGSurface) SendClose() {
	C.wlr_xdg_toplevel_send_close(s.p)
}

func (s XDGSurface) Ping() {
	C.wlr_xdg_surface_ping(s.p)
}

func (s XDGSurface) Surface() Surface {
	return Surface{p: s.p.surface}
}

func (s XDGSurface) SurfaceAt(sx float64, sy float64) (surface Surface, subX float64, subY float64) {
	var csubX, csubY C.double
	p := C.wlr_xdg_surface_surface_at(s.p, C.double(sx), C.double(sy), &csubX, &csubY)
	return Surface{p: p}, float64(csubX), float64(csubY)
}

func (s *XDGSurface) OnMap(cb func(*XDGSurface)) func() {
	lis := newListener(unsafe.Pointer(s.p), func(lis *wlrlis, data unsafe.Pointer) {
		cb(s)
	})
	C.wl_signal_add(&s.p.events._map, lis)
	return func() {
		removeListener(lis)
	}
}

func (s *XDGSurface) OnUnmap(cb func(*XDGSurface)) func() {
	lis := newListener(unsafe.Pointer(s.p), func(lis *wlrlis, data unsafe.Pointer) {
		cb(s)
	})
	C.wl_signal_add(&s.p.events.unmap, lis)
	return func() {
		removeListener(lis)
	}
}

func (s *XDGSurface) OnDestroy(cb func(*XDGSurface)) func() {
	lis := newListener(unsafe.Pointer(s.p), func(lis *wlrlis, data unsafe.Pointer) {
		cb(s)
	})
	C.wl_signal_add(&s.p.events.destroy, lis)
	return func() {
		removeListener(lis)
	}
}

func (s *XDGSurface) OnPingTimeout(cb func(*XDGSurface)) func() {
	lis := newListener(unsafe.Pointer(s.p), func(lis *wlrlis, data unsafe.Pointer) {
		cb(s)
	})
	C.wl_signal_add(&s.p.events.ping_timeout, lis)
	return func() {
		removeListener(lis)
	}
}

func (s *XDGSurface) OnNewPopup(cb func(*XDGSurface, *XDGPopup)) func() {
	lis := newListener(unsafe.Pointer(s.p), func(lis *wlrlis, data unsafe.Pointer) {
		popup := &XDGPopup{p: (*C.struct_wlr_xdg_popup)(data)}
		cb(s, popup)
	})
	C.wl_signal_add(&s.p.events.ping_timeout, lis)
	return func() {
		removeListener(lis)
	}
}

func (s XDGSurface) Geometry() Box {
	var cb C.struct_wlr_box
	C.wlr_xdg_surface_get_geometry(s.p, &cb)

	var b Box
	b.fromC(&cb)
	return b
}

type XDGPopup struct {
	p *C.struct_wlr_xdg_popup
}

type XDGTopLevel struct {
	p *C.struct_wlr_xdg_toplevel
}

func (t *XDGTopLevel) OnRequestMove(cb func(client *SeatClient, serial uint32)) func() {
	lis := newListener(unsafe.Pointer(t.p), func(lis *wlrlis, data unsafe.Pointer) {
		event := (*C.struct_wlr_xdg_toplevel_move_event)(data)
		client := &SeatClient{p: event.seat}
		cb(client, uint32(event.serial))
	})
	C.wl_signal_add(&t.p.events.request_move, lis)
	return func() {
		removeListener(lis)
	}
}

func (t *XDGTopLevel) OnRequestResize(cb func(client *SeatClient, serial uint32, edges Edges)) func() {
	lis := newListener(unsafe.Pointer(t.p), func(lis *wlrlis, data unsafe.Pointer) {
		event := (*C.struct_wlr_xdg_toplevel_resize_event)(data)
		client := &SeatClient{p: event.seat}
		cb(client, uint32(event.serial), Edges(event.edges))
	})
	C.wl_signal_add(&t.p.events.request_resize, lis)
	return func() {
		removeListener(lis)
	}
}

func (s XDGTopLevel) Nil() bool {
	return s.p == nil
}

func (t XDGTopLevel) Title() string {
	return C.GoString(t.p.title)
}

//export _wlr_xdg_surface_for_each_cb
func _wlr_xdg_surface_for_each_cb(surface *C.struct_wlr_surface, sx C.int, sy C.int, data unsafe.Pointer) {
	xdgSurfaceWalkersMutex.RLock()
	cb := xdgSurfaceWalkers[(*C.struct_wlr_xdg_surface)(data)]
	xdgSurfaceWalkersMutex.RUnlock()
	if cb != nil {
		cb(Surface{p: surface}, int(sx), int(sy))
	}
}
