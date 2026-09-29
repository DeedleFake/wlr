package wlr

/*
#include <wlr/xcursor.h>
#include <wlr/xwayland.h>
*/
import "C"

import "unsafe"

type XwaylandSurfaceDecorations uint32

const (
	XwaylandSurfaceDecorationsAll      XwaylandSurfaceDecorations = C.WLR_XWAYLAND_SURFACE_DECORATIONS_ALL
	XwaylandSurfaceDecorationsNoBorder XwaylandSurfaceDecorations = C.WLR_XWAYLAND_SURFACE_DECORATIONS_NO_BORDER
	XwaylandSurfaceDecorationsNoTitle  XwaylandSurfaceDecorations = C.WLR_XWAYLAND_SURFACE_DECORATIONS_NO_TITLE
)

type Xwayland struct {
	p *C.struct_wlr_xwayland
}

func CreateXwayland(display Display, compositor Compositor, lazy bool) Xwayland {
	p := C.wlr_xwayland_create(display.p, compositor.p, C.bool(lazy))
	return Xwayland{p: p}
}

func (x Xwayland) Valid() bool {
	return x.p != nil
}

func (x Xwayland) Destroy() {
	C.wlr_xwayland_destroy(x.p)
}

func (x Xwayland) Server() XwaylandServer {
	return XwaylandServer{p: x.p.server}
}

func (x Xwayland) OnNewSurface(cb func(XwaylandSurface)) Listener {
	return newListener(&x.p.events.new_surface, func(lis Listener, data unsafe.Pointer) {
		cb(XwaylandSurface{p: (*C.struct_wlr_xwayland_surface)(data)})
	})
}

func (x Xwayland) SetCursor(img XCursorImage) {
	buf := C.wlr_xcursor_image_get_buffer(img.p)
	C.wlr_xwayland_set_cursor(x.p, buf, C.int32_t(img.p.hotspot_x), C.int32_t(img.p.hotspot_y))
}

func (x Xwayland) DisplayName() string {
	return C.GoString(x.p.display_name)
}

type XwaylandSurface struct {
	p *C.struct_wlr_xwayland_surface
}

func (s XwaylandSurface) Valid() bool {
	return s.p != nil
}

func (s XwaylandSurface) Width() int {
	return int(s.p.width)
}

func (s XwaylandSurface) Height() int {
	return int(s.p.height)
}

func (s XwaylandSurface) Title() string {
	return C.GoString(s.p.title)
}

func (s XwaylandSurface) Decorations() XwaylandSurfaceDecorations {
	return XwaylandSurfaceDecorations(s.p.decorations)
}

func (s XwaylandSurface) Surface() Surface {
	return Surface{p: s.p.surface}
}

func (s XwaylandSurface) Close() {
	C.wlr_xwayland_surface_close(s.p)
}

func (s XwaylandSurface) Activate(a bool) {
	C.wlr_xwayland_surface_activate(s.p, C.bool(a))
}

func (s XwaylandSurface) SetMinimized(minimized bool) {
	C.wlr_xwayland_surface_set_minimized(s.p, C.bool(minimized))
}

func (s XwaylandSurface) SetMaximized(maximizedHorz, maximizedVert bool) {
	C.wlr_xwayland_surface_set_maximized(s.p, C.bool(maximizedHorz), C.bool(maximizedVert))
}

func (s XwaylandSurface) Configure(x int16, y int16, width uint16, height uint16) {
	C.wlr_xwayland_surface_configure(s.p, C.int16_t(x), C.int16_t(y), C.uint16_t(width), C.uint16_t(height))
}

func (s XwaylandSurface) OnDestroy(cb func(XwaylandSurface)) Listener {
	return newListener(&s.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

// OnAssociate is called when the X11 window is paired with a
// wlr_surface. Surface returns an invalid Surface until then, so
// listeners on the Surface, such as OnMap, should be added here.
func (s XwaylandSurface) OnAssociate(cb func(XwaylandSurface)) Listener {
	return newListener(&s.p.events.associate, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

// OnDissociate is called when the X11 window is about to lose its
// wlr_surface. Listeners added to the Surface in OnAssociate should be
// removed here.
func (s XwaylandSurface) OnDissociate(cb func(XwaylandSurface)) Listener {
	return newListener(&s.p.events.dissociate, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s XwaylandSurface) OnRequestMove(cb func(surface XwaylandSurface)) Listener {
	return newListener(&s.p.events.request_move, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s XwaylandSurface) OnRequestResize(cb func(surface XwaylandSurface, edges Edges)) Listener {
	return newListener(&s.p.events.request_resize, func(lis Listener, data unsafe.Pointer) {
		event := (*C.struct_wlr_xwayland_resize_event)(data)
		cb(s, Edges(event.edges))
	})
}

func (s XwaylandSurface) OnRequestMinimize(cb func(surface XwaylandSurface)) Listener {
	return newListener(&s.p.events.request_minimize, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s XwaylandSurface) OnRequestMaximize(cb func(surface XwaylandSurface)) Listener {
	return newListener(&s.p.events.request_maximize, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s XwaylandSurface) OnRequestConfigure(cb func(surface XwaylandSurface, x int16, y int16, width uint16, height uint16)) Listener {
	return newListener(&s.p.events.request_configure, func(lis Listener, data unsafe.Pointer) {
		event := (*C.struct_wlr_xwayland_surface_configure_event)(data)
		cb(s, int16(event.x), int16(event.y), uint16(event.width), uint16(event.height))
	})
}

func (s XwaylandSurface) OnSetTitle(cb func(XwaylandSurface, string)) Listener {
	return newListener(&s.p.events.set_title, func(lis Listener, data unsafe.Pointer) {
		cb(s, C.GoString((*C.char)(data)))
	})
}

func (s XwaylandSurface) OnSetDecorations(cb func(XwaylandSurface)) Listener {
	return newListener(&s.p.events.set_decorations, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

type XwaylandServer struct {
	p *C.struct_wlr_xwayland_server
}

func (s XwaylandServer) DisplayName() string {
	return C.GoString(&s.p.display_name[0])
}
