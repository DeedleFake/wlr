package wlr

/*
#include <stdlib.h>
#include <wayland-server-core.h>

struct _listener_data {
	struct wl_listener lis;
	uintptr_t handle;
};

static inline struct _listener_data *_listener_get_data(struct wl_listener *lis) {
	struct _listener_data *data;
	return wl_container_of(lis, data, lis);
}

extern void _listener_callback(struct wl_listener *listener, void *data);

static inline void _listener_set_callback(struct wl_listener *listener) {
	listener->notify = _listener_callback;
}
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

type listenerFunc func(lis Listener, data unsafe.Pointer)

// Listener represents an attached signal handler for a Wayland event
// of some kind.
//
// Note: It is the client's responsibility to call Destroy when they
// are done with a Listener in order to free resources. Failure to do
// so will result in a memory leak.
type Listener struct {
	p *C.struct__listener_data
}

func newListener(sig *C.struct_wl_signal, cb listenerFunc) Listener {
	lis := Listener{
		p: (*C.struct__listener_data)(C.malloc(C.sizeof_struct__listener_data)),
	}
	lis.p.handle = C.uintptr_t(cgo.NewHandle(cb))
	C._listener_set_callback(&lis.p.lis)

	if sig != nil {
		C.wl_signal_add(sig, &lis.p.lis)
	}

	return lis
}

// Destroy frees resources associated with the Listener and disconnects
// it from the signal that it is attached to. After this is called,
// Valid will return false.
//
// Calling this method on an invalid Listener is a no-op.
func (lis Listener) Destroy() {
	if !lis.Valid() {
		return
	}

	cgo.Handle(lis.p.handle).Delete()
	C.wl_list_remove(&lis.p.lis.link)
	C.free(unsafe.Pointer(lis.p))
	lis.p = nil
}

func (lis Listener) Valid() bool {
	return lis.p != nil
}

// call calls the Listener with the given data. What exactly the data
// is used for depends on how the Listener was created.
func (lis Listener) call(data unsafe.Pointer) {
	cgo.Handle(lis.p.handle).Value().(listenerFunc)(lis, data)
}

//export _listener_callback
func _listener_callback(p *C.struct_wl_listener, data unsafe.Pointer) {
	lis := Listener{p: C._listener_get_data(p)}
	lis.call(data)
}
