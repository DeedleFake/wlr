package wlr

// #include <wlr/types/wlr_input_device.h>
// #include <wlr/types/wlr_keyboard.h>
// #include <wlr/types/wlr_pointer.h>
import "C"

import (
	"fmt"
	"unsafe"

	"deedles.dev/wlr/xkb"
)

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
