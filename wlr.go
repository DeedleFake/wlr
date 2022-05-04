package wlr

// #cgo pkg-config: wlroots wayland-server
// #cgo CFLAGS: -D_GNU_SOURCE -DWLR_USE_UNSTABLE
//
// #include <stdarg.h>
// #include <stdio.h>
// #include <stdlib.h>
// #include <time.h>
// #include <wayland-server.h>
//
// #include <wlr/backend/wayland.h>
// #include <wlr/backend/x11.h>
// #include <wlr/util/box.h>
// #include <wlr/types/wlr_compositor.h>
// #include <wlr/types/wlr_cursor.h>
// #include <wlr/types/wlr_data_device.h>
// #include <wlr/types/wlr_server_decoration.h>
// #include <wlr/types/wlr_linux_dmabuf_v1.h>
// #include <wlr/types/wlr_input_device.h>
// #include <wlr/types/wlr_keyboard.h>
// #include <wlr/types/wlr_matrix.h>
// #include <wlr/types/wlr_output.h>
// #include <wlr/types/wlr_output_layout.h>
// #include <wlr/types/wlr_seat.h>
// #include <wlr/types/wlr_xcursor_manager.h>
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
	"errors"
	"fmt"
	"sync"
	"time"
	"unsafe"

	"deedles.dev/wlr/xkb"
)

type XCursor struct {
	p *C.struct_wlr_xcursor
}

type XCursorImage struct {
	p *C.struct_wlr_xcursor_image
}

type XCursorManager struct {
	p *C.struct_wlr_xcursor_manager
}

func NewXCursorManager() XCursorManager {
	p := C.wlr_xcursor_manager_create(nil, 24)
	return XCursorManager{p: p}
}

func (m XCursorManager) Destroy() {
	C.wlr_xcursor_manager_destroy(m.p)
}

func (m XCursorManager) Load() {
	C.wlr_xcursor_manager_load(m.p, 1)
}

func (m XCursorManager) SetCursorImage(cursor Cursor, name string) {
	s := C.CString(name)
	C.wlr_xcursor_manager_set_cursor_image(m.p, s, cursor.p)
	C.free(unsafe.Pointer(s))
}

func (m XCursorManager) XCursor(name string, scale float32) XCursor {
	s := C.CString(name)
	p := C.wlr_xcursor_manager_get_xcursor(m.p, s, C.float(scale))
	C.free(unsafe.Pointer(s))
	return XCursor{p: p}
}

func (c XCursor) Image(i int) XCursorImage {
	n := c.ImageCount()
	slice := (*[1 << 30]*C.struct_wlr_xcursor_image)(unsafe.Pointer(c.p.images))[:n:n]
	return XCursorImage{p: slice[i]}
}

func (c XCursor) Images() []XCursorImage {
	images := make([]XCursorImage, 0, c.ImageCount())
	for i := 0; i < cap(images); i++ {
		images = append(images, c.Image(i))
	}
	return images
}

func (c XCursor) ImageCount() int {
	return int(c.p.image_count)
}

func (c XCursor) Name() string {
	return C.GoString(c.p.name)
}

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

type Seat struct {
	p *C.struct_wlr_seat
}

type SeatClient struct {
	p *C.struct_wlr_seat_client
}

type SeatKeyboardState struct {
	s C.struct_wlr_seat_keyboard_state
}

type SeatPointerState struct {
	s C.struct_wlr_seat_pointer_state
}

type SeatCapability uint32

const (
	SeatCapabilityPointer  SeatCapability = C.WL_SEAT_CAPABILITY_POINTER
	SeatCapabilityKeyboard SeatCapability = C.WL_SEAT_CAPABILITY_KEYBOARD
	SeatCapabilityTouch    SeatCapability = C.WL_SEAT_CAPABILITY_TOUCH
)

func NewSeat(display Display, name string) Seat {
	s := C.CString(name)
	p := C.wlr_seat_create(display.p, s)
	C.free(unsafe.Pointer(s))
	man.track(unsafe.Pointer(p), &p.events.destroy)
	return Seat{p: p}
}

func (s Seat) Destroy() {
	C.wlr_seat_destroy(s.p)
}

func (s Seat) OnDestroy(cb func(Seat)) {
	man.add(unsafe.Pointer(s.p), &s.p.events.destroy, func(unsafe.Pointer) {
		cb(s)
	})
}

func (s Seat) OnSetCursorRequest(cb func(client SeatClient, surface Surface, serial uint32, hotspotX int32, hotspotY int32)) {
	man.add(unsafe.Pointer(s.p), &s.p.events.request_set_cursor, func(data unsafe.Pointer) {
		event := (*C.struct_wlr_seat_pointer_request_set_cursor_event)(data)
		client := SeatClient{p: event.seat_client}
		surface := Surface{p: event.surface}
		cb(client, surface, uint32(event.serial), int32(event.hotspot_x), int32(event.hotspot_y))
	})
}

func (s Seat) SetCapabilities(caps SeatCapability) {
	C.wlr_seat_set_capabilities(s.p, C.uint32_t(caps))
}

func (s Seat) SetKeyboard(dev InputDevice) {
	C.wlr_seat_set_keyboard(s.p, dev.p)
}

func (s Seat) NotifyPointerButton(time uint32, button uint32, state ButtonState) {
	C.wlr_seat_pointer_notify_button(s.p, C.uint32_t(time), C.uint32_t(button), uint32(state))
}

func (s Seat) NotifyPointerAxis(time uint32, orientation AxisOrientation, delta float64, deltaDiscrete int32, source AxisSource) {
	C.wlr_seat_pointer_notify_axis(s.p, C.uint32_t(time), C.enum_wlr_axis_orientation(orientation), C.double(delta), C.int32_t(deltaDiscrete), C.enum_wlr_axis_source(source))
}

func (s Seat) NotifyPointerEnter(surface Surface, sx float64, sy float64) {
	C.wlr_seat_pointer_notify_enter(s.p, surface.p, C.double(sx), C.double(sy))
}

func (s Seat) NotifyPointerMotion(time uint32, sx float64, sy float64) {
	C.wlr_seat_pointer_notify_motion(s.p, C.uint32_t(time), C.double(sx), C.double(sy))
}

func (s Seat) NotifyPointerFrame() {
	C.wlr_seat_pointer_notify_frame(s.p)
}

func (s Seat) NotifyKeyboardEnter(surface Surface, k Keyboard) {
	C.wlr_seat_keyboard_notify_enter(s.p, surface.p, &k.p.keycodes[0], k.p.num_keycodes, &k.p.modifiers)
}

func (s Seat) NotifyKeyboardModifiers(k Keyboard) {
	C.wlr_seat_keyboard_notify_modifiers(s.p, &k.p.modifiers)
}

func (s Seat) NotifyKeyboardKey(time uint32, keyCode uint32, state KeyState) {
	C.wlr_seat_keyboard_notify_key(s.p, C.uint32_t(time), C.uint32_t(keyCode), C.uint32_t(state))
}

func (s Seat) ClearPointerFocus() {
	C.wlr_seat_pointer_clear_focus(s.p)
}

func (s Seat) Keyboard() Keyboard {
	p := C.wlr_seat_get_keyboard(s.p)
	return Keyboard{p: p}
}

func (s Seat) KeyboardState() SeatKeyboardState {
	return SeatKeyboardState{s: s.p.keyboard_state}
}

func (s Seat) PointerState() SeatPointerState {
	return SeatPointerState{s: s.p.pointer_state}
}

func (s SeatKeyboardState) FocusedSurface() Surface {
	return Surface{p: s.s.focused_surface}
}

func (s SeatPointerState) FocusedSurface() Surface {
	return Surface{p: s.s.focused_surface}
}

func (s SeatPointerState) FocusedClient() SeatClient {
	return SeatClient{p: s.s.focused_client}
}

type Output struct {
	p *C.struct_wlr_output
}

type OutputMode struct {
	p *C.struct_wlr_output_mode
}

func wrapOutput(p unsafe.Pointer) Output {
	return Output{p: (*C.struct_wlr_output)(p)}
}

func (o Output) OnDestroy(cb func(Output)) {
	man.add(unsafe.Pointer(o.p), &o.p.events.destroy, func(unsafe.Pointer) {
		cb(o)
	})
}

func (o Output) Name() string {
	return C.GoString(o.p.name)
}

func (o Output) Scale() float32 {
	return float32(o.p.scale)
}

func (o Output) TransformMatrix() Matrix {
	var matrix Matrix
	matrix.fromC(&o.p.transform_matrix)
	return matrix
}

func (o Output) OnFrame(cb func(Output)) {
	man.add(unsafe.Pointer(o.p), &o.p.events.frame, func(data unsafe.Pointer) {
		cb(o)
	})
}

func (o Output) RenderSoftwareCursors() {
	C.wlr_output_render_software_cursors(o.p, nil)
}

func (o Output) TransformedResolution() (int, int) {
	var width, height C.int
	C.wlr_output_transformed_resolution(o.p, &width, &height)
	return int(width), int(height)
}

func (o Output) EffectiveResolution() (int, int) {
	var width, height C.int
	C.wlr_output_effective_resolution(o.p, &width, &height)
	return int(width), int(height)
}

func (o Output) AttachRender() (int, error) {
	var bufferAge C.int
	if !C.wlr_output_attach_render(o.p, &bufferAge) {
		return 0, errors.New("can't make output context current")
	}

	return int(bufferAge), nil
}

func (o Output) Rollback() {
	C.wlr_output_rollback(o.p)
}

func (o Output) CreateGlobal() {
	C.wlr_output_create_global(o.p)
}

func (o Output) DestroyGlobal() {
	C.wlr_output_destroy_global(o.p)
}

func (o Output) Commit() {
	C.wlr_output_commit(o.p)
}

func (o Output) Modes() []OutputMode {
	// TODO: figure out what to do with this ridiculous for loop
	// perhaps this can be refactored into a less ugly hack that uses reflection
	var modes []OutputMode
	var mode *C.struct_wlr_output_mode
	for mode := (*C.struct_wlr_output_mode)(unsafe.Pointer(uintptr(unsafe.Pointer(o.p.modes.next)) - unsafe.Offsetof(mode.link))); &mode.link != &o.p.modes; mode = (*C.struct_wlr_output_mode)(unsafe.Pointer(uintptr(unsafe.Pointer(mode.link.next)) - unsafe.Offsetof(mode.link))) {
		modes = append(modes, OutputMode{p: mode})
	}

	return modes
}

func (o Output) SetMode(mode OutputMode) {
	C.wlr_output_set_mode(o.p, mode.p)
}

func (o Output) Enable(enable bool) {
	C.wlr_output_enable(o.p, C.bool(enable))
}

func (o Output) SetTitle(title string) error {
	if C.wlr_output_is_wl(o.p) {
		C.wlr_wl_output_set_title(o.p, C.CString(title))
	} else if C.wlr_output_is_x11(o.p) {
		C.wlr_x11_output_set_title(o.p, C.CString(title))
	} else {
		return errors.New("this output type cannot have a title")
	}

	return nil
}

type OutputLayout struct {
	p *C.struct_wlr_output_layout
}

func NewOutputLayout() OutputLayout {
	p := C.wlr_output_layout_create()
	man.track(unsafe.Pointer(p), &p.events.destroy)
	return OutputLayout{p: p}
}

func (l OutputLayout) Destroy() {
	C.wlr_output_layout_destroy(l.p)
}

func (l OutputLayout) AddOutputAuto(output Output) {
	C.wlr_output_layout_add_auto(l.p, output.p)
}

func (l OutputLayout) Coords(output Output) (x float64, y float64) {
	var ox, oy C.double
	C.wlr_output_layout_output_coords(l.p, output.p, &ox, &oy)
	return float64(ox), float64(oy)
}

func OutputTransformInvert(transform uint32) uint32 {
	return uint32(C.wlr_output_transform_invert(C.enum_wl_output_transform(transform)))
}

func (m OutputMode) Width() int32 {
	return int32(m.p.width)
}

func (m OutputMode) Height() int32 {
	return int32(m.p.height)
}

func (m OutputMode) RefreshRate() int32 {
	return int32(m.p.refresh)
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

type (
	KeyState         uint32
	KeyboardModifier uint32
)

const (
	KeyStateReleased KeyState = C.WL_KEYBOARD_KEY_STATE_RELEASED
	KeyStatePressed  KeyState = C.WL_KEYBOARD_KEY_STATE_PRESSED

	KeyboardModifierShift KeyboardModifier = C.WLR_MODIFIER_SHIFT
	KeyboardModifierCaps  KeyboardModifier = C.WLR_MODIFIER_CAPS
	KeyboardModifierCtrl  KeyboardModifier = C.WLR_MODIFIER_CTRL
	KeyboardModifierAlt   KeyboardModifier = C.WLR_MODIFIER_ALT
	KeyboardModifierMod2  KeyboardModifier = C.WLR_MODIFIER_MOD2
	KeyboardModifierMod3  KeyboardModifier = C.WLR_MODIFIER_MOD3
	KeyboardModifierLogo  KeyboardModifier = C.WLR_MODIFIER_LOGO
	KeyboardModifierMod5  KeyboardModifier = C.WLR_MODIFIER_MOD5
)

type Keyboard struct {
	p *C.struct_wlr_keyboard
}

func (k Keyboard) SetKeymap(keymap xkb.Keymap) {
	C.wlr_keyboard_set_keymap(k.p, (*C.struct_xkb_keymap)(keymap.Ptr()))
}

func (k Keyboard) RepeatInfo() (rate int32, delay int32) {
	return int32(k.p.repeat_info.rate), int32(k.p.repeat_info.delay)
}

func (k Keyboard) SetRepeatInfo(rate int32, delay int32) {
	C.wlr_keyboard_set_repeat_info(k.p, C.int32_t(rate), C.int32_t(delay))
}

func (k Keyboard) XKBState() xkb.State {
	return xkb.WrapState(unsafe.Pointer(k.p.xkb_state))
}

func (k Keyboard) Modifiers() KeyboardModifier {
	return KeyboardModifier(C.wlr_keyboard_get_modifiers(k.p))
}

func (k Keyboard) OnModifiers(cb func(keyboard Keyboard)) {
	man.add(unsafe.Pointer(k.p), &k.p.events.modifiers, func(data unsafe.Pointer) {
		cb(k)
	})
}

func (k Keyboard) OnKey(cb func(keyboard Keyboard, time uint32, keyCode uint32, updateState bool, state KeyState)) {
	man.add(unsafe.Pointer(k.p), &k.p.events.key, func(data unsafe.Pointer) {
		event := (*C.struct_wlr_event_keyboard_key)(data)
		cb(k, uint32(event.time_msec), uint32(event.keycode), bool(event.update_state), KeyState(event.state))
	})
}

type (
	InputDeviceType uint32
	ButtonState     uint32
	AxisSource      uint32
	AxisOrientation uint32
)

var inputDeviceNames = []string{
	InputDeviceTypeKeyboard:   "keyboard",
	InputDeviceTypePointer:    "pointer",
	InputDeviceTypeTouch:      "touch",
	InputDeviceTypeTabletTool: "tablet tool",
	InputDeviceTypeTabletPad:  "tablet pad",
}

const (
	InputDeviceTypeKeyboard   InputDeviceType = C.WLR_INPUT_DEVICE_KEYBOARD
	InputDeviceTypePointer    InputDeviceType = C.WLR_INPUT_DEVICE_POINTER
	InputDeviceTypeTouch      InputDeviceType = C.WLR_INPUT_DEVICE_TOUCH
	InputDeviceTypeTabletTool InputDeviceType = C.WLR_INPUT_DEVICE_TABLET_TOOL
	InputDeviceTypeTabletPad  InputDeviceType = C.WLR_INPUT_DEVICE_TABLET_PAD

	ButtonStateReleased ButtonState = C.WLR_BUTTON_RELEASED
	ButtonStatePressed  ButtonState = C.WLR_BUTTON_PRESSED

	AxisSourceWheel      AxisSource = C.WLR_AXIS_SOURCE_WHEEL
	AxisSourceFinger     AxisSource = C.WLR_AXIS_SOURCE_FINGER
	AxisSourceContinuous AxisSource = C.WLR_AXIS_SOURCE_CONTINUOUS
	AxisSourceWheelTilt  AxisSource = C.WLR_AXIS_SOURCE_WHEEL_TILT

	AxisOrientationVertical   AxisOrientation = C.WLR_AXIS_ORIENTATION_VERTICAL
	AxisOrientationHorizontal AxisOrientation = C.WLR_AXIS_ORIENTATION_HORIZONTAL
)

type InputDevice struct {
	p *C.struct_wlr_input_device
}

func (d InputDevice) OnDestroy(cb func(InputDevice)) {
	man.add(unsafe.Pointer(d.p), &d.p.events.destroy, func(unsafe.Pointer) {
		cb(d)
	})
}

func (d InputDevice) Type() InputDeviceType { return InputDeviceType(d.p._type) }
func (d InputDevice) Vendor() int           { return int(d.p.vendor) }
func (d InputDevice) Product() int          { return int(d.p.product) }
func (d InputDevice) Name() string          { return C.GoString(d.p.name) }
func (d InputDevice) Width() float64        { return float64(d.p.width_mm) }
func (d InputDevice) Height() float64       { return float64(d.p.height_mm) }
func (d InputDevice) OutputName() string    { return C.GoString(d.p.output_name) }

func validateInputDeviceType(d InputDevice, fn string, req InputDeviceType) {
	if typ := d.Type(); typ != req {
		if int(typ) >= len(inputDeviceNames) {
			panic(fmt.Sprintf("%s called on input device of type %d", fn, typ))
		} else {
			panic(fmt.Sprintf("%s called on input device of type %s", fn, inputDeviceNames[typ]))
		}
	}
}

func (d InputDevice) Keyboard() Keyboard {
	validateInputDeviceType(d, "Keyboard", InputDeviceTypeKeyboard)
	p := *(*unsafe.Pointer)(unsafe.Pointer(&d.p.anon0[0]))
	return Keyboard{p: (*C.struct_wlr_keyboard)(p)}
}

func wrapInputDevice(p unsafe.Pointer) InputDevice {
	return InputDevice{p: (*C.struct_wlr_input_device)(p)}
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

type Display struct {
	p *C.struct_wl_display
}

func NewDisplay() Display {
	p := C.wl_display_create()
	d := Display{p: p}
	d.OnDestroy(func(Display) {
		man.delete(unsafe.Pointer(p))
	})
	return d
}

func (d Display) Destroy() {
	C.wl_display_destroy(d.p)
}

func (d Display) OnDestroy(cb func(Display)) {
	l := man.add(unsafe.Pointer(d.p), nil, func(data unsafe.Pointer) {
		cb(d)
	})
	C.wl_display_add_destroy_listener(d.p, l.p)
}

func (d Display) Run() {
	C.wl_display_run(d.p)
}

func (d Display) Terminate() {
	C.wl_display_terminate(d.p)
}

func (d Display) EventLoop() EventLoop {
	p := C.wl_display_get_event_loop(d.p)
	evl := EventLoop{p: p}
	evl.OnDestroy(func(EventLoop) {
		man.delete(unsafe.Pointer(p))
	})
	return evl
}

func (d Display) AddSocketAuto() (string, error) {
	socket := C.wl_display_add_socket_auto(d.p)
	if socket == nil {
		return "", errors.New("can't auto add wayland socket")
	}

	return C.GoString(socket), nil
}

func (d Display) FlushClients() {
	C.wl_display_flush_clients(d.p)
}

type ServerDecorationManagerMode uint32

const (
	ServerDecorationManagerModeNone   ServerDecorationManagerMode = C.WLR_SERVER_DECORATION_MANAGER_MODE_NONE
	ServerDecorationManagerModeClient ServerDecorationManagerMode = C.WLR_SERVER_DECORATION_MANAGER_MODE_CLIENT
	ServerDecorationManagerModeServer ServerDecorationManagerMode = C.WLR_SERVER_DECORATION_MANAGER_MODE_SERVER
)

type ServerDecorationManager struct {
	p *C.struct_wlr_server_decoration_manager
}

type ServerDecoration struct {
	p *C.struct_wlr_server_decoration
}

func NewServerDecorationManager(display Display) ServerDecorationManager {
	p := C.wlr_server_decoration_manager_create(display.p)
	man.track(unsafe.Pointer(p), &p.events.destroy)
	return ServerDecorationManager{p: p}
}

func (m ServerDecorationManager) OnDestroy(cb func(ServerDecorationManager)) {
	man.add(unsafe.Pointer(m.p), &m.p.events.destroy, func(unsafe.Pointer) {
		cb(m)
	})
}

func (m ServerDecorationManager) SetDefaultMode(mode ServerDecorationManagerMode) {
	C.wlr_server_decoration_manager_set_default_mode(m.p, C.uint32_t(mode))
}

func (m ServerDecorationManager) OnNewMode(cb func(ServerDecorationManager, ServerDecoration)) {
	man.add(unsafe.Pointer(m.p), &m.p.events.new_decoration, func(data unsafe.Pointer) {
		dec := ServerDecoration{
			p: (*C.struct_wlr_server_decoration)(data),
		}
		man.track(unsafe.Pointer(dec.p), &dec.p.events.destroy)
		cb(m, dec)
	})
}

func (d ServerDecoration) OnDestroy(cb func(ServerDecoration)) {
	man.add(unsafe.Pointer(d.p), &d.p.events.destroy, func(unsafe.Pointer) {
		cb(d)
	})
}

func (d ServerDecoration) OnMode(cb func(ServerDecoration)) {
	man.add(unsafe.Pointer(d.p), &d.p.events.mode, func(unsafe.Pointer) {
		cb(d)
	})
}

func (d ServerDecoration) Mode() ServerDecorationManagerMode {
	return ServerDecorationManagerMode(d.p.mode)
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

type Cursor struct {
	p *C.struct_wlr_cursor
}

func NewCursor() Cursor {
	p := C.wlr_cursor_create()
	return Cursor{p: p}
}

func (c Cursor) Destroy() {
	C.wlr_cursor_destroy(c.p)
	man.delete(unsafe.Pointer(c.p))
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

func (c Cursor) OnMotion(cb func(dev InputDevice, time uint32, dx float64, dy float64)) {
	man.add(unsafe.Pointer(c.p), &c.p.events.motion, func(data unsafe.Pointer) {
		event := (*C.struct_wlr_event_pointer_motion)(data)
		dev := InputDevice{p: event.device}
		cb(dev, uint32(event.time_msec), float64(event.delta_x), float64(event.delta_y))
	})
}

func (c Cursor) OnMotionAbsolute(cb func(dev InputDevice, time uint32, x float64, y float64)) {
	man.add(unsafe.Pointer(c.p), &c.p.events.motion_absolute, func(data unsafe.Pointer) {
		event := (*C.struct_wlr_event_pointer_motion_absolute)(data)
		dev := InputDevice{p: event.device}
		cb(dev, uint32(event.time_msec), float64(event.x), float64(event.y))
	})
}

func (c Cursor) OnButton(cb func(dev InputDevice, time uint32, button uint32, state ButtonState)) {
	man.add(unsafe.Pointer(c.p), &c.p.events.button, func(data unsafe.Pointer) {
		event := (*C.struct_wlr_event_pointer_button)(data)
		dev := InputDevice{p: event.device}
		cb(dev, uint32(event.time_msec), uint32(event.button), ButtonState(event.state))
	})
}

func (c Cursor) OnAxis(cb func(dev InputDevice, time uint32, source AxisSource, orientation AxisOrientation, delta float64, deltaDiscrete int32)) {
	man.add(unsafe.Pointer(c.p), &c.p.events.axis, func(data unsafe.Pointer) {
		event := (*C.struct_wlr_event_pointer_axis)(data)
		dev := InputDevice{p: event.device}
		cb(dev, uint32(event.time_msec), AxisSource(event.source), AxisOrientation(event.orientation), float64(event.delta), int32(event.delta_discrete))
	})
}

func (c Cursor) OnFrame(cb func()) {
	man.add(unsafe.Pointer(c.p), &c.p.events.frame, func(data unsafe.Pointer) {
		cb()
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