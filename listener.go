package wlr

/*
#include <stdlib.h>
#include <wayland-server-core.h>

extern void _listener_callback(struct wl_listener *listener, void *data);

static inline void _listener_set_callback(struct wl_listener *listener) {
	listener->notify = _listener_callback;
}

static inline void *_listener_get_handle(struct wl_listener *listener) {
	return listener + 1;
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
	p *C.struct_wl_listener
}

func newListener(sig *C.struct_wl_signal, cb listenerFunc) Listener {
	// I guess that wl_listener can do user data after all. Huh.
	lis := Listener{
		p: (*C.struct_wl_listener)(C.malloc(C.sizeof_struct_wl_listener + C.sizeof_uintptr_t)),
	}
	*(lis.handle()) = cgo.NewHandle(cb)
	C._listener_set_callback(lis.p)

	if sig != nil {
		C.wl_signal_add(sig, lis.p)
	}

	return lis
}

func (lis Listener) handle() *cgo.Handle {
	return (*cgo.Handle)(C._listener_get_handle(lis.p))
}

// Destroy frees resources associated with the Listener and disconnects
// it from the signal that it is attached to.
//
// The behavior of this method if called twice is undefined.
func (lis Listener) Destroy() {
	lis.handle().Delete()
	C.wl_list_remove(&lis.p.link)
	C.free(unsafe.Pointer(lis.p))
	lis.p = nil
}

func (lis Listener) Valid() bool {
	return lis.p != nil
}

//export _listener_callback
func _listener_callback(p *C.struct_wl_listener, data unsafe.Pointer) {
	lis := Listener{p: p}
	lis.handle().Value().(listenerFunc)(lis, data)
}
