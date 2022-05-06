package wlr

// #include <stdlib.h>
// #include <wayland-server-core.h>
//
// extern void _listener_callback(struct wl_listener *listener, void *data);
//
// static inline void _set_listener_callback(struct wl_listener *listener) {
// 	listener->notify = _listener_callback;
// }
import "C"
import (
	"fmt"
	"unsafe"

	"deedles.dev/wlr/internal/util"
)

type wlrlis = C.struct_wl_listener

var cbs util.SMap[*wlrlis, callback]

type callback struct {
	obj unsafe.Pointer
	cb  callbackFunc
}

type callbackFunc func(lis *wlrlis, data unsafe.Pointer)

func newListener(obj unsafe.Pointer, cb callbackFunc) *wlrlis {
	lis := (*wlrlis)(C.malloc(C.sizeof_struct_wl_listener))
	C._set_listener_callback(lis)
	cbs.Store(lis, callback{obj, cb})
	return lis
}

func removeListener(lis *wlrlis) {
	cbs.Delete(lis)
	C.wl_list_remove(&lis.link)
	C.free(unsafe.Pointer(lis))
}

func trackObject(p unsafe.Pointer, sig *C.struct_wl_signal) func() {
	lis := newListener(p, func(lis *wlrlis, data unsafe.Pointer) {
		removeObject(unsafe.Pointer(p))
	})
	C.wl_signal_add(sig, lis)
	return func() {
		removeListener(lis)
	}
}

func removeObject(obj unsafe.Pointer) {
	cbs.Range(func(lis *wlrlis, cb callback) bool {
		if cb.obj != obj {
			return true
		}
		removeListener(lis)
		return true
	})
}

//export _listener_callback
func _listener_callback(lis *wlrlis, data unsafe.Pointer) {
	cb, ok := cbs.Load(lis)
	if !ok {
		panic(fmt.Errorf("no callback found for listener %v with data %v", lis, data))
	}
	cb.cb(lis, data)
}
