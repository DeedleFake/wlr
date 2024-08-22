package wlr

/*
#include <wlr/types/wlr_compositor.h>
#include <wlr/types/wlr_subcompositor.h>
#include <wlr/types/wlr_xdg_shell.h>
#include <wlr/xwayland.h>

void _wlr_surface_for_each_cb(struct wlr_surface *surface, int sx, int sy, void *data);

static inline void _wlr_surface_for_each_surface(struct wlr_surface *surface, void *user_data) {
	wlr_surface_for_each_surface(surface, _wlr_surface_for_each_cb, user_data);
}

struct _wlr_surface_has_surface_data {
	struct wlr_surface *sub;
	int found;
};

static void _wlr_surface_has_surface_cb(struct wlr_surface *surface, int sx, int sy, void *d) {
	struct _wlr_surface_has_surface_data *data = (struct _wlr_surface_has_surface_data *)d;
	if (surface == data->sub) {
		data->found = 1;
	}
}

static inline int _wlr_surface_has_surface(struct wlr_surface *surface, struct wlr_surface *sub) {
	struct _wlr_surface_has_surface_data data = (struct _wlr_surface_has_surface_data){
		.sub = sub,
		.found = 0,
	};
	wlr_surface_for_each_surface(surface, _wlr_surface_has_surface_cb, &data);
	return data.found;
}

// _new_timespec exists to avoid possible problems from differing type
// names on some systems. C is less picky than Go, so it shouldn't be
// a problem if it's done here.
static inline struct timespec _new_timespec(long sec, long nsec) {
	return (struct timespec){
		.tv_sec = sec,
		.tv_nsec = nsec,
	};
}
*/
import "C"

import (
	"runtime/cgo"
	"time"
	"unsafe"
)

type Surface struct {
	p *C.struct_wlr_surface
}

func (s Surface) Valid() bool {
	return s.p != nil
}

func (s Surface) OnMap(cb func(Surface)) Listener {
	return newListener(&s.p.events._map, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s Surface) OnUnmap(cb func(Surface)) Listener {
	return newListener(&s.p.events.unmap, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s Surface) OnDestroy(cb func(Surface)) Listener {
	return newListener(&s.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

func (s Surface) Mapped() bool {
	return bool(s.p.mapped)
}

func (s Surface) SurfaceAt(sx float64, sy float64) (surface Surface, subX, subY float64, ok bool) {
	var csubX, csubY C.double
	p := C.wlr_surface_surface_at(s.p, C.double(sx), C.double(sy), &csubX, &csubY)
	return Surface{p: p}, float64(csubX), float64(csubY), p != nil
}

func (s Surface) GetTexture() Texture {
	p := C.wlr_surface_get_texture(s.p)
	return Texture{p: p}
}

func (s Surface) Current() SurfaceState {
	return SurfaceState{v: s.p.current}
}

func (s Surface) ForEachSurface(cb func(Surface, int, int)) {
	handle := cgo.NewHandle(cb)
	defer handle.Delete()

	C._wlr_surface_for_each_surface(s.p, unsafe.Pointer(&handle))
}

// HasSurface is a convenience function that searches for sub in s. It
// does the search entirely in C, so it may be more effecient than
// manually iterating.
func (s Surface) HasSurface(sub Surface) bool {
	return C._wlr_surface_has_surface(s.p, sub.p) != 0
}

//export _wlr_surface_for_each_cb
func _wlr_surface_for_each_cb(surface *C.struct_wlr_surface, sx C.int, sy C.int, data unsafe.Pointer) {
	handle := *(*cgo.Handle)(data)
	cb := handle.Value().(func(Surface, int, int))
	cb(Surface{p: surface}, int(sx), int(sy))
}

func (s Surface) SendEnter(output Output) {
	C.wlr_surface_send_enter(s.p, output.p)
}

func (s Surface) SendLeave(output Output) {
	C.wlr_surface_send_leave(s.p, output.p)
}

func (s Surface) SendFrameDone(when time.Time) {
	ts := C._new_timespec(C.long(when.Unix()), C.long(when.Nanosecond()))
	C.wlr_surface_send_frame_done(s.p, &ts)
}

func (s Surface) XwaylandSurface() XwaylandSurface {
	p := C.wlr_xwayland_surface_try_from_wlr_surface(s.p)
	return XwaylandSurface{p: p}
}

func (s Surface) XDGSurface() XDGSurface {
	p := C.wlr_xdg_surface_try_from_wlr_surface(s.p)
	return XDGSurface{p: p}
}

type SurfaceState struct {
	v C.struct_wlr_surface_state
}

func (s SurfaceState) Dx() int32 {
	return int32(s.v.dx)
}

func (s SurfaceState) Dy() int32 {
	return int32(s.v.dy)
}

func (s SurfaceState) Width() int {
	return int(s.v.width)
}

func (s SurfaceState) Height() int {
	return int(s.v.height)
}

func (s SurfaceState) Transform() OutputTransform {
	return OutputTransform(s.v.transform)
}
