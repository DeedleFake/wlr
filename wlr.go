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

	"golang.org/x/exp/constraints"
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

func matrixFromC(cm *[9]C.float) *Matrix {
	var m Matrix
	copyMatrix((*[9]float32)(&m), cm)
	return &m
}

func copyMatrix[
	Out, In constraints.Integer | constraints.Float,
](out *[9]Out, in *[9]In) {
	for i := range out {
		out[i] = Out(in[i])
	}
}

func ProjectBoxMatrix(box image.Rectangle, transform OutputTransform, rotation float32, projection *Matrix) *Matrix {
	var cm [9]C.float
	pm := projection.toC()
	C.wlr_matrix_project_box(
		&cm[0],
		boxToC(box),
		C.enum_wl_output_transform(transform),
		C.float(rotation),
		&pm[0],
	)
	return matrixFromC(&cm)
}

func (m *Matrix) toC() [9]C.float {
	var cm [9]C.float
	copyMatrix(&cm, (*[9]float32)(m))
	return cm
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

func CreateCompositor(display Display, version uint32, renderer Renderer) Compositor {
	p := C.wlr_compositor_create(display.p, C.uint32_t(version), renderer.p)
	return Compositor{p: p}
}

func (c Compositor) OnDestroy(cb func(Compositor)) Listener {
	return newListener(&c.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(c)
	})
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

func boxToC(r image.Rectangle) *C.struct_wlr_box {
	if r == image.ZR {
		return nil
	}

	r = r.Canon()
	return &C.struct_wlr_box{
		x:      C.int(r.Min.X),
		y:      C.int(r.Min.Y),
		width:  C.int(r.Dx()),
		height: C.int(r.Dy()),
	}
}

func boxFromC(box *C.struct_wlr_box) image.Rectangle {
	if box == nil {
		return image.ZR
	}

	return image.Rect(
		int(box.x),
		int(box.y),
		int(box.x+box.width),
		int(box.y+box.height),
	)
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

type Resource struct {
	p *C.struct_wl_resource
}

func (r Resource) GetClient() Client {
	return Client{p: C.wl_resource_get_client(r.p)}
}

type Client struct {
	p *C.struct_wl_client
}

func (c Client) GetCredentials() (pid, uid, gid int) {
	var cpid C.pid_t
	var cuid C.uid_t
	var cgid C.gid_t
	C.wl_client_get_credentials(c.p, &cpid, &cuid, &cgid)
	return int(cpid), int(cuid), int(cgid)
}

func container[T any](list *C.struct_wl_list, offset int) *T {
	return (*T)(unsafe.Add(unsafe.Pointer(list), -offset))
}

func listForEach[T any](head *C.struct_wl_list, offset int, do func(*T) bool) {
	pos := head.next
	for {
		if head == pos {
			return
		}

		if !do(container[T](pos, offset)) {
			return
		}

		pos = pos.next
	}
}
