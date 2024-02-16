package wlr

/*
#include <stdlib.h>
#include <linux/input-event-codes.h>

#include <wayland-server-core.h>
#include <wlr/types/wlr_cursor.h>
#include <wlr/types/wlr_pointer.h>
*/
import "C"

import (
	"time"
	"unsafe"
)

type Cursor struct {
	p *C.struct_wlr_cursor
}

func CreateCursor() Cursor {
	p := C.wlr_cursor_create()
	return Cursor{p: p}
}

func (c Cursor) Destroy() {
	C.wlr_cursor_destroy(c.p)
}

func (c Cursor) X() float64 {
	return float64(c.p.x)
}

func (c Cursor) Y() float64 {
	return float64(c.p.y)
}

func (c Cursor) AttachOutputLayout(layout OutputLayout) {
	C.wlr_cursor_attach_output_layout(c.p, layout.p)
}

func (c Cursor) AttachInputDevice(dev InputDevice) {
	C.wlr_cursor_attach_input_device(c.p, dev.p)
}

func (c Cursor) Move(dev InputDevice, dx float64, dy float64) {
	C.wlr_cursor_move(c.p, dev.p, C.double(dx), C.double(dy))
}

func (c Cursor) WarpAbsolute(dev InputDevice, x float64, y float64) {
	C.wlr_cursor_warp_absolute(c.p, dev.p, C.double(x), C.double(y))
}

func (c Cursor) SetSurface(surface Surface, hotspotX int32, hotspotY int32) {
	C.wlr_cursor_set_surface(c.p, surface.p, C.int32_t(hotspotX), C.int32_t(hotspotY))
}

func (c Cursor) SetXCursor(m XCursorManager, name string) {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))

	C.wlr_cursor_set_xcursor(c.p, m.p, s)
}

func (c Cursor) OnMotion(cb func(p Pointer, time time.Time, dx, dy float64)) Listener {
	return newListener(&c.p.events.motion, func(lis Listener, data unsafe.Pointer) {
		event := (*C.struct_wlr_pointer_motion_event)(data)
		dev := Pointer{p: event.pointer}
		cb(dev, time.UnixMilli(int64(event.time_msec)), float64(event.delta_x), float64(event.delta_y))
	})
}

func (c Cursor) OnMotionAbsolute(cb func(p Pointer, time time.Time, x, y float64)) Listener {
	return newListener(&c.p.events.motion_absolute, func(lis Listener, data unsafe.Pointer) {
		event := (*C.struct_wlr_pointer_motion_absolute_event)(data)
		dev := Pointer{p: event.pointer}
		cb(dev, time.UnixMilli(int64(event.time_msec)), float64(event.x), float64(event.y))
	})
}

func (c Cursor) OnButton(cb func(p Pointer, time time.Time, button CursorButton, state ButtonState)) Listener {
	return newListener(&c.p.events.button, func(lis Listener, data unsafe.Pointer) {
		event := (*C.struct_wlr_pointer_button_event)(data)
		dev := Pointer{p: event.pointer}
		cb(dev, time.UnixMilli(int64(event.time_msec)), CursorButton(event.button), ButtonState(event.state))
	})
}

func (c Cursor) OnAxis(cb func(p Pointer, time time.Time, source AxisSource, orientation AxisOrientation, delta float64, deltaDiscrete int32)) Listener {
	return newListener(&c.p.events.axis, func(lis Listener, data unsafe.Pointer) {
		event := (*C.struct_wlr_pointer_axis_event)(data)
		dev := Pointer{p: event.pointer}
		cb(
			dev,
			time.UnixMilli(int64(event.time_msec)),
			AxisSource(event.source),
			AxisOrientation(event.orientation),
			float64(event.delta),
			int32(event.delta_discrete),
		)
	})
}

func (c Cursor) OnFrame(cb func()) Listener {
	return newListener(&c.p.events.frame, func(lis Listener, data unsafe.Pointer) {
		cb()
	})
}

type CursorButton uint32

const (
	BtnLeft   CursorButton = C.BTN_LEFT
	BtnRight  CursorButton = C.BTN_RIGHT
	BtnMiddle CursorButton = C.BTN_MIDDLE
)

type Pointer struct {
	p *C.struct_wlr_pointer
}

func (p Pointer) Base() InputDevice {
	return InputDevice{p: &p.p.base}
}

func (p Pointer) OutputName() string {
	return C.GoString(p.p.output_name)
}
