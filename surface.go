package wlr

/*
#include <wlr/types/wlr_surface.h>
#include <wlr/types/wlr_xdg_shell.h>
#include <wlr/xwayland.h>

void _wlr_surface_for_each_cb(struct wlr_surface *surface, int sx, int sy, void *data);

static inline void _wlr_surface_for_each_surface(struct wlr_surface *surface, void *user_data) {
	wlr_surface_for_each_surface(surface, _wlr_surface_for_each_cb, user_data);
}
*/
import "C"

import (
	"runtime/cgo"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

type SurfaceType uint32

const (
	SurfaceTypeNone SurfaceType = iota
	SurfaceTypeXDG
	SurfaceTypeXWayland
)

type Surface struct {
	p *C.struct_wlr_surface
}

func (s Surface) Valid() bool {
	return s.p != nil
}

func (s Surface) OnDestroy(cb func(Surface)) Listener {
	return newListener(&s.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s Surface) Type() SurfaceType {
	if C.wlr_surface_is_xdg_surface(s.p) {
		return SurfaceTypeXDG
	} else if C.wlr_surface_is_xwayland_surface(s.p) {
		return SurfaceTypeXWayland
	}

	return SurfaceTypeNone
}

func (s Surface) SurfaceAt(sx float64, sy float64) (surface Surface, subX float64, subY float64) {
	var csubX, csubY C.double
	p := C.wlr_surface_surface_at(s.p, C.double(sx), C.double(sy), &csubX, &csubY)
	return Surface{p: p}, float64(csubX), float64(csubY)
}

func (s Surface) Texture() Texture {
	p := C.wlr_surface_get_texture(s.p)
	return Texture{p: p}
}

func (s Surface) Current() SurfaceState {
	return SurfaceState{s: s.p.current}
}

func (s Surface) Pending() SurfaceState {
	return SurfaceState{s: s.p.pending}
}

func (s Surface) ForEachSurface(cb func(Surface, int, int)) {
	handle := cgo.NewHandle(cb)
	C._wlr_surface_for_each_surface(s.p, unsafe.Pointer(&handle))
}

//export _wlr_surface_for_each_cb
func _wlr_surface_for_each_cb(surface *C.struct_wlr_surface, sx C.int, sy C.int, data unsafe.Pointer) {
	handle := *(*cgo.Handle)(data)
	defer handle.Delete()

	cb := handle.Value().(func(Surface, int, int))
	cb(Surface{p: surface}, int(sx), int(sy))
}

func (s Surface) SendFrameDone(when time.Time) {
	// we ignore the returned error; the only possible error is
	// ERANGE, when timespec on a platform has int32 precision, but
	// our time requires 64 bits. This should not occur.
	t, _ := unix.TimeToTimespec(when)
	C.wlr_surface_send_frame_done(s.p, (*C.struct_timespec)(unsafe.Pointer(&t)))
}

func (s Surface) XDGSurface() XDGSurface {
	p := C.wlr_xdg_surface_from_wlr_surface(s.p)
	return XDGSurface{p: p}
}

func (s Surface) XWaylandSurface() XWaylandSurface {
	p := C.wlr_xwayland_surface_from_wlr_surface(s.p)
	return XWaylandSurface{p: p}
}

type SurfaceState struct {
	s C.struct_wlr_surface_state
}

func (s SurfaceState) Width() int {
	return int(s.s.width)
}

func (s SurfaceState) Height() int {
	return int(s.s.height)
}

func (s SurfaceState) Transform() uint32 {
	return uint32(s.s.transform)
}
