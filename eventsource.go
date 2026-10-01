package wlr

/*
#include <stdint.h>
#include <wayland-server-core.h>

extern int _event_source_fd_callback(int fd, uint32_t mask, void *data);

static inline struct wl_event_source *_event_loop_add_fd(struct wl_event_loop *loop, int fd, uint32_t mask, uintptr_t handle) {
	return wl_event_loop_add_fd(loop, fd, mask, _event_source_fd_callback, (void *)handle);
}
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

type EventMask uint32

const (
	EventReadable EventMask = C.WL_EVENT_READABLE
	EventWritable EventMask = C.WL_EVENT_WRITABLE
	EventHangup   EventMask = C.WL_EVENT_HANGUP
	EventError    EventMask = C.WL_EVENT_ERROR
)

// EventSource is a source of events attached to an EventLoop.
//
// Note: It is the client's responsibility to call Remove when they are
// done with an EventSource in order to free resources.
type EventSource struct {
	p *C.struct_wl_event_source
	h cgo.Handle
}

type eventSourceFdFunc func(fd uintptr, mask EventMask)

// AddFd calls cb on the event loop whenever fd is ready for any of the
// events in mask. EventHangup and EventError are always reported.
func (evl EventLoop) AddFd(fd uintptr, mask EventMask, cb func(fd uintptr, mask EventMask)) EventSource {
	h := cgo.NewHandle(eventSourceFdFunc(cb))
	p := C._event_loop_add_fd(evl.p, C.int(fd), C.uint32_t(mask), C.uintptr_t(h))
	if p == nil {
		h.Delete()
		return EventSource{}
	}
	return EventSource{p: p, h: h}
}

func (s EventSource) Valid() bool {
	return s.p != nil
}

// Remove removes the source from its event loop and frees it.
func (s EventSource) Remove() {
	C.wl_event_source_remove(s.p)
	s.h.Delete()
}

//export _event_source_fd_callback
func _event_source_fd_callback(fd C.int, mask C.uint32_t, data unsafe.Pointer) C.int {
	cgo.Handle(uintptr(data)).Value().(eventSourceFdFunc)(uintptr(fd), EventMask(mask))
	return 0
}
