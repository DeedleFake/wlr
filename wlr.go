package wlr

/*
#cgo pkg-config: wlroots wayland-server pixman-1
#cgo CFLAGS: -D_GNU_SOURCE -DWLR_USE_UNSTABLE

#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

#include <wlr/util/box.h>
#include <wlr/types/wlr_compositor.h>
#include <wlr/types/wlr_data_device.h>
#include <wlr/types/wlr_matrix.h>
#include <wlr/util/edges.h>
#include <wlr/xwayland.h>
*/
import "C"

import (
	"image"
	"image/color"
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

type EventLoop struct {
	p *C.struct_wl_event_loop
}

func (evl EventLoop) OnDestroy(cb func(EventLoop)) Listener {
	lis := newListener(nil, func(lis Listener, data unsafe.Pointer) {
		cb(evl)
	})
	C.wl_event_loop_add_destroy_listener(evl.p, &lis.p.lis)
	return lis
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

func CreateDataDeviceManager(display Display) DataDeviceManager {
	p := C.wlr_data_device_manager_create(display.p)
	return DataDeviceManager{p: p}
}

func (m DataDeviceManager) OnDestroy(cb func(DataDeviceManager)) Listener {
	return newListener(&m.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(m)
	})
}

type Compositor struct {
	p *C.struct_wlr_compositor
}

func CreateCompositor(display Display, renderer Renderer) Compositor {
	p := C.wlr_compositor_create(display.p, renderer.p)
	return Compositor{p: p}
}

func (c Compositor) OnDestroy(cb func(Compositor)) Listener {
	return newListener(&c.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(c)
	})
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

func colorToC(c color.Color) [4]C.float {
	r, g, b, a := c.RGBA()
	return [...]C.float{
		C.float(r) / C.float(a),
		C.float(g) / C.float(a),
		C.float(b) / C.float(a),
		C.float(a) / 0xFFFF,
	}
}

func rectToC(r image.Rectangle) (cr C.pixman_region32_t) {
	r = r.Canon()
	C.pixman_region32_init_rect(
		&cr,
		C.int(r.Min.X),
		C.int(r.Max.Y),
		C.uint(r.Dx()),
		C.uint(r.Dy()),
	)
	return
}
