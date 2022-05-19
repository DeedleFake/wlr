package wlr

/*
#include <wlr/xwayland.h>
*/
import "C"

import "unsafe"

type XWayland struct {
	p *C.struct_wlr_xwayland
}

func CreateXWayland(display Display, compositor Compositor, lazy bool) XWayland {
	p := C.wlr_xwayland_create(display.p, compositor.p, C.bool(lazy))
	return XWayland{p: p}
}

func (x XWayland) Valid() bool {
	return x.p != nil
}

func (x XWayland) Destroy() {
	C.wlr_xwayland_destroy(x.p)
}

func (x XWayland) Server() XWaylandServer {
	return XWaylandServer{p: x.p.server}
}

func (x XWayland) OnNewSurface(cb func(XWaylandSurface)) Listener {
	return newListener(&x.p.events.new_surface, func(lis Listener, data unsafe.Pointer) {
		cb(XWaylandSurface{p: (*C.struct_wlr_xwayland_surface)(data)})
	})
}

func (x XWayland) SetCursor(img XCursorImage) {
	C.wlr_xwayland_set_cursor(x.p, img.p.buffer, img.p.width*4, img.p.width, img.p.height, C.int32_t(img.p.hotspot_x), C.int32_t(img.p.hotspot_y))
}

type XWaylandSurface struct {
	p *C.struct_wlr_xwayland_surface
}

func (s XWaylandSurface) Valid() bool {
	return s.p != nil
}

func (s XWaylandSurface) Mapped() bool {
	return bool(s.p.mapped)
}

func (s XWaylandSurface) Width() int {
	return int(s.p.width)
}

func (s XWaylandSurface) Height() int {
	return int(s.p.height)
}

func (s XWaylandSurface) Title() string {
	return C.GoString(s.p.title)
}

func (s XWaylandSurface) Surface() Surface {
	return Surface{p: s.p.surface}
}

func (s XWaylandSurface) Close() {
	C.wlr_xwayland_surface_close(s.p)
}

func (s XWaylandSurface) Activate(a bool) {
	C.wlr_xwayland_surface_activate(s.p, C.bool(a))
}

func (s XWaylandSurface) SetMinimized(minimized bool) {
	C.wlr_xwayland_surface_set_minimized(s.p, C.bool(minimized))
}

func (s XWaylandSurface) SetMaximized(maximized bool) {
	C.wlr_xwayland_surface_set_maximized(s.p, C.bool(maximized))
}

func (s XWaylandSurface) Configure(x int16, y int16, width uint16, height uint16) {
	C.wlr_xwayland_surface_configure(s.p, C.int16_t(x), C.int16_t(y), C.uint16_t(width), C.uint16_t(height))
}

func (s XWaylandSurface) OnMap(cb func(XWaylandSurface)) Listener {
	return newListener(&s.p.events._map, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s XWaylandSurface) OnUnmap(cb func(XWaylandSurface)) Listener {
	return newListener(&s.p.events.unmap, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s XWaylandSurface) OnDestroy(cb func(XWaylandSurface)) Listener {
	return newListener(&s.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s XWaylandSurface) OnRequestMove(cb func(surface XWaylandSurface)) Listener {
	return newListener(&s.p.events.request_move, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s XWaylandSurface) OnRequestResize(cb func(surface XWaylandSurface, edges Edges)) Listener {
	return newListener(&s.p.events.request_resize, func(lis Listener, data unsafe.Pointer) {
		event := (*C.struct_wlr_xwayland_resize_event)(data)
		cb(s, Edges(event.edges))
	})
}

func (s XWaylandSurface) OnRequestMinimize(cb func(surface XWaylandSurface)) Listener {
	return newListener(&s.p.events.request_minimize, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s XWaylandSurface) OnRequestMaximize(cb func(surface XWaylandSurface)) Listener {
	return newListener(&s.p.events.request_maximize, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s XWaylandSurface) OnRequestConfigure(cb func(surface XWaylandSurface, x int16, y int16, width uint16, height uint16)) Listener {
	return newListener(&s.p.events.request_configure, func(lis Listener, data unsafe.Pointer) {
		event := (*C.struct_wlr_xwayland_surface_configure_event)(data)
		cb(s, int16(event.x), int16(event.y), uint16(event.width), uint16(event.height))
	})
}

func (s XWaylandSurface) OnSetTitle(cb func(XWaylandSurface, string)) Listener {
	return newListener(&s.p.events.set_title, func(lis Listener, data unsafe.Pointer) {
		cb(s, C.GoString((*C.char)(data)))
	})
}

type XWaylandServer struct {
	p *C.struct_wlr_xwayland_server
}

func (s XWaylandServer) DisplayName() string {
	return C.GoString(&s.p.display_name[0])
}
