package wlr

// #cgo pkg-config: wlroots wayland-server
// #cgo CFLAGS: -D_GNU_SOURCE -DWLR_USE_UNSTABLE
//
// #include <stdarg.h>
// #include <stdio.h>
// #include <stdlib.h>
// #include <time.h>
//
// #include <wlr/util/box.h>
// #include <wlr/types/wlr_compositor.h>
// #include <wlr/types/wlr_data_device.h>
// #include <wlr/types/wlr_linux_dmabuf_v1.h>
// #include <wlr/types/wlr_matrix.h>
// #include <wlr/util/edges.h>
// #include <wlr/util/log.h>
// #include <wlr/xwayland.h>
//
// void _wlr_log_cb(enum wlr_log_importance importance, char *msg);
//
// static inline void _wlr_log_inner_cb(enum wlr_log_importance importance, const char *fmt, va_list args) {
// 	char *msg = NULL;
// 	if (vasprintf(&msg, fmt, args) == -1) {
// 		return;
// 	}
//
// 	_wlr_log_cb(importance, msg);
// 	free(msg);
// }
//
// static inline void _wlr_log_set_cb(enum wlr_log_importance verbosity, bool is_set) {
// 	wlr_log_init(verbosity, is_set ? &_wlr_log_inner_cb : NULL);
// }
//
// void _wl_listener_cb(struct wl_listener *listener, void *data);
//
// static inline void _wl_listener_set_cb(struct wl_listener *listener) {
// 	listener->notify = &_wl_listener_cb;
// }
import "C"

import (
	"sync"
	"time"
	"unsafe"
)

type Edges uint32

const (
	EdgeNone   Edges = C.WLR_EDGE_NONE
	EdgeTop    Edges = C.WLR_EDGE_TOP
	EdgeBottom Edges = C.WLR_EDGE_BOTTOM
	EdgeLeft   Edges = C.WLR_EDGE_LEFT
	EdgeRight  Edges = C.WLR_EDGE_RIGHT
)

type Texture struct {
	p *C.struct_wlr_texture
}

func (t Texture) Destroy() {
	C.wlr_texture_destroy(t.p)
}

func (t Texture) Nil() bool {
	return t.p == nil
}

type Matrix [9]float32

func (m *Matrix) ProjectBox(box *Box, transform uint32, rotation float32, projection *Matrix) {
	cm := m.toC()
	b := box.toC()
	pm := projection.toC()
	C.wlr_matrix_project_box(&cm[0], &b, C.enum_wl_output_transform(transform), C.float(rotation), &pm[0])
	m.fromC(&cm)
}

func (m *Matrix) toC() [9]C.float {
	var cm [9]C.float
	for i := range m {
		cm[i] = C.float(m[i])
	}
	return cm
}

func (m *Matrix) fromC(cm *[9]C.float) {
	for i := range cm {
		m[i] = float32(cm[i])
	}
}

type (
	LogImportance uint32
	LogFunc       func(importance LogImportance, msg string)
)

const (
	LogImportanceSilent LogImportance = C.WLR_SILENT
	LogImportanceError  LogImportance = C.WLR_ERROR
	LogImportanceInfo   LogImportance = C.WLR_INFO
	LogImportanceDebug  LogImportance = C.WLR_DEBUG
)

var (
	onLog LogFunc
)

//export _wlr_log_cb
func _wlr_log_cb(importance LogImportance, msg *C.char) {
	if onLog != nil {
		onLog(importance, C.GoString(msg))
	}
}

func OnLog(verbosity LogImportance, cb LogFunc) {
	C._wlr_log_set_cb(C.enum_wlr_log_importance(verbosity), cb != nil)
	onLog = cb
}

type DMABuf struct {
	p *C.struct_wlr_linux_dmabuf_v1
}

func NewDMABuf(display Display, renderer Renderer) DMABuf {
	p := C.wlr_linux_dmabuf_v1_create(display.p, renderer.p)
	man.track(unsafe.Pointer(p), &p.events.destroy)
	return DMABuf{p: p}
}

func (b DMABuf) OnDestroy(cb func(DMABuf)) {
	man.add(unsafe.Pointer(b.p), &b.p.events.destroy, func(unsafe.Pointer) {
		cb(b)
	})
}

type EventLoop struct {
	p *C.struct_wl_event_loop
}

func (evl EventLoop) OnDestroy(cb func(EventLoop)) {
	l := man.add(unsafe.Pointer(evl.p), nil, func(data unsafe.Pointer) {
		cb(evl)
	})
	C.wl_event_loop_add_destroy_listener(evl.p, l.p)
}

func (evl EventLoop) Fd() uintptr {
	return uintptr(C.wl_event_loop_get_fd(evl.p))
}

func (evl EventLoop) Dispatch(timeout time.Duration) {
	var d int
	if timeout >= 0 {
		d = int(timeout / time.Millisecond)
	} else {
		d = -1
	}
	C.wl_event_loop_dispatch(evl.p, C.int(d))
}

type DataDeviceManager struct {
	p *C.struct_wlr_data_device_manager
}

func NewDataDeviceManager(display Display) DataDeviceManager {
	p := C.wlr_data_device_manager_create(display.p)
	man.track(unsafe.Pointer(p), &p.events.destroy)
	return DataDeviceManager{p: p}
}

func (m DataDeviceManager) OnDestroy(cb func(DataDeviceManager)) {
	man.add(unsafe.Pointer(m.p), &m.p.events.destroy, func(unsafe.Pointer) {
		cb(m)
	})
}

type Compositor struct {
	p *C.struct_wlr_compositor
}

func NewCompositor(display Display, renderer Renderer) Compositor {
	p := C.wlr_compositor_create(display.p, renderer.p)
	man.track(unsafe.Pointer(p), &p.events.destroy)
	return Compositor{p: p}
}

func (c Compositor) OnDestroy(cb func(Compositor)) {
	man.add(unsafe.Pointer(c.p), &c.p.events.destroy, func(unsafe.Pointer) {
		cb(c)
	})
}

type Color struct {
	R, G, B, A float32
}

func (c *Color) Set(r, g, b, a float32) {
	c.R = r
	c.G = g
	c.B = b
	c.A = a
}

func (c *Color) toC() [4]C.float {
	return [...]C.float{
		C.float(c.R),
		C.float(c.G),
		C.float(c.B),
		C.float(c.A),
	}
}

type Box struct {
	X, Y, Width, Height int
}

func (b *Box) Set(x, y, width, height int) {
	b.X = x
	b.Y = y
	b.Width = width
	b.Height = height
}

func (b *Box) toC() C.struct_wlr_box {
	return C.struct_wlr_box{
		x:      C.int(b.X),
		y:      C.int(b.Y),
		width:  C.int(b.Width),
		height: C.int(b.Height),
	}
}

func (b *Box) fromC(cb *C.struct_wlr_box) {
	b.X = int(cb.x)
	b.Y = int(cb.y)
	b.Width = int(cb.width)
	b.Height = int(cb.height)
}

// This whole mess has to exist for a number of reasons:
//
// 1. We need to allocate all instances of wl_listener on the heap as storing Go
// pointers in C after a cgo call returns is not allowed.
//
// 2. The wlroots library implicitly destroys objects when wl_display is
// destroyed. So, we need to keep track of all objects (and their listeners)
// manually and listen for the destroy signal to be able to free everything.
//
// 3 (TODO). As we're keeping track of all objects anyway, we might as well
// store a Go pointer to the wrapper struct along with them in order to be able
// to pass the same Go pointer through callbacks every time. This will also
// allow calling runtime.SetFinalizer on some of them to clean them up early
// when the GC notices it has gone out of scope.
//
// Send help.

type (
	listenerCallback func(data unsafe.Pointer)
)

type manager struct {
	mutex     sync.RWMutex
	objects   map[unsafe.Pointer][]*listener
	listeners map[*C.struct_wl_listener]*listener
}

type listener struct {
	p   *C.struct_wl_listener
	s   *C.struct_wl_signal
	cbs []listenerCallback
}

var (
	man = &manager{
		objects:   map[unsafe.Pointer][]*listener{},
		listeners: map[*C.struct_wl_listener]*listener{},
	}
)

//export _wl_listener_cb
func _wl_listener_cb(listener *C.struct_wl_listener, data unsafe.Pointer) {
	man.mutex.RLock()
	l := man.listeners[listener]
	man.mutex.RUnlock()
	for _, cb := range l.cbs {
		cb(data)
	}
}

func (m *manager) add(p unsafe.Pointer, signal *C.struct_wl_signal, cb listenerCallback) *listener {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// if a listener for this object and signal already exists, add the callback
	// to the existing listener
	if signal != nil {
		for _, l := range m.objects[p] {
			if l.s != nil && l.s == signal {
				l.cbs = append(l.cbs, cb)
				return l
			}
		}
	}

	lp := (*C.struct_wl_listener)(C.calloc(C.sizeof_struct_wl_listener, 1))
	C._wl_listener_set_cb(lp)
	if signal != nil {
		C.wl_signal_add((*C.struct_wl_signal)(signal), lp)
	}

	l := &listener{
		p:   lp,
		s:   signal,
		cbs: []listenerCallback{cb},
	}
	m.listeners[lp] = l
	m.objects[p] = append(m.objects[p], l)

	return l
}

func (m *manager) has(p unsafe.Pointer) bool {
	m.mutex.RLock()
	_, found := m.objects[p]
	m.mutex.RUnlock()
	return found
}

func (m *manager) delete(p unsafe.Pointer) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, l := range m.objects[p] {
		delete(m.listeners, l.p)

		// remove the listener from the signal
		C.wl_list_remove(&l.p.link)

		// free the listener
		C.free(unsafe.Pointer(l.p))
	}

	delete(m.objects, p)
}

func (m *manager) track(p unsafe.Pointer, destroySignal *C.struct_wl_signal) {
	m.add(p, destroySignal, func(data unsafe.Pointer) { m.delete(p) })
}
