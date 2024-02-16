package wlr

/*
#include <wlr/backend.h>
#include <wlr/types/wlr_xdg_shell.h>
#include <wlr/types/wlr_xdg_output_v1.h>

void _wlr_surface_for_each_cb(struct wlr_surface *surface, int sx, int sy, void *data);

static inline void _wlr_xdg_surface_for_each_surface(struct wlr_xdg_surface *surface, void *user_data) {
	wlr_xdg_surface_for_each_surface(surface, _wlr_surface_for_each_cb, user_data);
}

struct _wlr_xdg_surface_has_surface_data {
	struct wlr_surface *sub;
	int found;
};

static void _wlr_xdg_surface_has_surface_cb(struct wlr_surface *surface, int sx, int sy, void *d) {
	struct _wlr_xdg_surface_has_surface_data *data = (struct _wlr_xdg_surface_has_surface_data *)d;
	if (surface == data->sub) {
		data->found = 1;
	}
}

static inline int _wlr_xdg_surface_has_surface(struct wlr_xdg_surface *surface, struct wlr_surface *sub) {
	struct _wlr_xdg_surface_has_surface_data data = (struct _wlr_xdg_surface_has_surface_data){
		.sub = sub,
		.found = 0,
	};
	wlr_xdg_surface_for_each_surface(surface, _wlr_xdg_surface_has_surface_cb, &data);
	return data.found;
}
*/
import "C"

import (
	"image"
	"runtime/cgo"
	"unsafe"
)

type XDGSurfaceRole uint32

const (
	XDGSurfaceRoleNone     XDGSurfaceRole = C.WLR_XDG_SURFACE_ROLE_NONE
	XDGSurfaceRoleToplevel XDGSurfaceRole = C.WLR_XDG_SURFACE_ROLE_TOPLEVEL
	XDGSurfaceRolePopup    XDGSurfaceRole = C.WLR_XDG_SURFACE_ROLE_POPUP
)

type XDGSurfaceWalkFunc func(surface Surface, sx int, sy int)

type XDGShell struct {
	p *C.struct_wlr_xdg_shell
}

func CreateXDGShell(display Display, version uint32) XDGShell {
	p := C.wlr_xdg_shell_create(display.p, C.uint32_t(version))
	return XDGShell{p: p}
}

func (s XDGShell) OnDestroy(cb func(XDGShell)) Listener {
	return newListener(&s.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s XDGShell) OnNewSurface(cb func(XDGSurface)) Listener {
	return newListener(&s.p.events.new_surface, func(lis Listener, data unsafe.Pointer) {
		cb(XDGSurface{p: (*C.struct_wlr_xdg_surface)(data)})
	})
}

type XDGSurface struct {
	p *C.struct_wlr_xdg_surface
}

func (s XDGSurface) Valid() bool {
	return s.p != nil
}

func (s XDGSurface) Resource() Resource {
	return Resource{p: s.p.resource}
}

func (s XDGSurface) ForEachSurface(cb func(Surface, int, int)) {
	handle := cgo.NewHandle(cb)
	defer handle.Delete()

	C._wlr_xdg_surface_for_each_surface(s.p, unsafe.Pointer(&handle))
}

func (s XDGSurface) HasSurface(sub Surface) bool {
	return C._wlr_xdg_surface_has_surface(s.p, sub.p) != 0
}

func (s XDGSurface) Role() XDGSurfaceRole {
	return XDGSurfaceRole(s.p.role)
}

func (s XDGSurface) Current() XDGSurfaceState {
	return XDGSurfaceState{v: s.p.current}
}

func (s XDGSurface) Toplevel() XDGToplevel {
	p := *(*unsafe.Pointer)(unsafe.Pointer(&s.p.anon0[0]))
	return XDGToplevel{p: (*C.struct_wlr_xdg_toplevel)(p)}
}

func (s XDGSurface) Popup() XDGPopup {
	p := *(*unsafe.Pointer)(unsafe.Pointer(&s.p.anon0[0]))
	return XDGPopup{p: (*C.struct_wlr_xdg_popup)(p)}
}

func (s XDGToplevel) SetActivated(activated bool) {
	C.wlr_xdg_toplevel_set_activated(s.p, C.bool(activated))
}

func (s XDGToplevel) SetResizing(resizing bool) {
	C.wlr_xdg_toplevel_set_resizing(s.p, C.bool(resizing))
}

func (s XDGToplevel) SetSize(width int32, height int32) {
	C.wlr_xdg_toplevel_set_size(s.p, C.int32_t(width), C.int32_t(height))
}

func (s XDGToplevel) SetTiled(edges Edges) {
	C.wlr_xdg_toplevel_set_tiled(s.p, C.uint32_t(edges))
}

func (s XDGToplevel) SetMaximized(maximized bool) {
	C.wlr_xdg_toplevel_set_maximized(s.p, C.bool(maximized))
}

func (s XDGToplevel) SendClose() {
	C.wlr_xdg_toplevel_send_close(s.p)
}

func (s XDGSurface) Ping() {
	C.wlr_xdg_surface_ping(s.p)
}

func (s XDGSurface) Surface() Surface {
	return Surface{p: s.p.surface}
}

func (s XDGSurface) SurfaceAt(sx float64, sy float64) (surface Surface, subX float64, subY float64, ok bool) {
	var csubX, csubY C.double
	p := C.wlr_xdg_surface_surface_at(s.p, C.double(sx), C.double(sy), &csubX, &csubY)
	return Surface{p: p}, float64(csubX), float64(csubY), p != nil
}

func (s XDGSurface) OnDestroy(cb func(XDGSurface)) Listener {
	return newListener(&s.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s XDGSurface) OnPingTimeout(cb func(XDGSurface)) Listener {
	return newListener(&s.p.events.ping_timeout, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s XDGSurface) OnNewPopup(cb func(XDGSurface, XDGPopup)) Listener {
	return newListener(&s.p.events.new_popup, func(lis Listener, data unsafe.Pointer) {
		cb(
			s,
			XDGPopup{p: (*C.struct_wlr_xdg_popup)(data)},
		)
	})
}

func (s XDGSurface) GetGeometry() image.Rectangle {
	var cb C.struct_wlr_box
	C.wlr_xdg_surface_get_geometry(s.p, &cb)
	return boxFromC(&cb)
}

type XDGPopup struct {
	p *C.struct_wlr_xdg_popup
}

func (p XDGPopup) Parent() Surface {
	return Surface{p: p.p.parent}
}

type XDGToplevel struct {
	p *C.struct_wlr_xdg_toplevel
}

func (t XDGToplevel) OnRequestMove(cb func(t XDGToplevel, client SeatClient, serial uint32)) Listener {
	return newListener(&t.p.events.request_move, func(lis Listener, data unsafe.Pointer) {
		event := (*C.struct_wlr_xdg_toplevel_move_event)(data)
		client := SeatClient{p: event.seat}
		cb(t, client, uint32(event.serial))
	})
}

func (t XDGToplevel) OnRequestResize(cb func(t XDGToplevel, client SeatClient, serial uint32, edges Edges)) Listener {
	return newListener(&t.p.events.request_resize, func(lis Listener, data unsafe.Pointer) {
		event := (*C.struct_wlr_xdg_toplevel_resize_event)(data)
		client := SeatClient{p: event.seat}
		cb(t, client, uint32(event.serial), Edges(event.edges))
	})
}

func (t XDGToplevel) OnRequestMinimize(cb func(XDGToplevel)) Listener {
	return newListener(&t.p.events.request_minimize, func(lis Listener, data unsafe.Pointer) {
		cb(t)
	})
}

func (t XDGToplevel) OnRequestMaximize(cb func(XDGToplevel)) Listener {
	return newListener(&t.p.events.request_maximize, func(lis Listener, data unsafe.Pointer) {
		cb(t)
	})
}

func (t XDGToplevel) OnSetTitle(cb func(XDGToplevel, string)) Listener {
	return newListener(&t.p.events.set_title, func(lis Listener, data unsafe.Pointer) {
		cb(t, C.GoString((*C.char)(data)))
	})
}

func (s XDGToplevel) Valid() bool {
	return s.p != nil
}

func (t XDGToplevel) Title() string {
	return C.GoString(t.p.title)
}

func (t XDGToplevel) Current() XDGToplevelState {
	return XDGToplevelState{v: t.p.current}
}

type XDGToplevelState struct {
	v C.struct_wlr_xdg_toplevel_state
}

func (s XDGToplevelState) Activated() bool {
	return bool(s.v.activated)
}

func (s XDGToplevelState) Width() uint32 {
	return uint32(s.v.width)
}

func (s XDGToplevelState) Height() uint32 {
	return uint32(s.v.height)
}

func (s XDGToplevelState) MinWidth() uint32 {
	return uint32(s.v.min_width)
}

func (s XDGToplevelState) MinHeight() uint32 {
	return uint32(s.v.min_height)
}

func (s XDGToplevelState) MaxWidth() uint32 {
	return uint32(s.v.max_width)
}

func (s XDGToplevelState) MaxHeight() uint32 {
	return uint32(s.v.max_height)
}

type XDGOutputManagerV1 struct {
	p *C.struct_wlr_xdg_output_manager_v1
}

func CreateXDGOutputManagerV1(display Display, layout OutputLayout) XDGOutputManagerV1 {
	p := C.wlr_xdg_output_manager_v1_create(display.p, layout.p)
	return XDGOutputManagerV1{p: p}
}

func (m XDGOutputManagerV1) OnDestroy(cb func(XDGOutputManagerV1)) Listener {
	return newListener(&m.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(m)
	})
}

type XDGSurfaceState struct {
	v C.struct_wlr_xdg_surface_state
}

func (s XDGSurfaceState) Geometry() image.Rectangle {
	return boxFromC(&s.v.geometry)
}
