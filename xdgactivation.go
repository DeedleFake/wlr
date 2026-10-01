package wlr

/*
#include <wlr/types/wlr_xdg_activation_v1.h>
*/
import "C"

import "unsafe"

type XDGActivationV1 struct {
	p *C.struct_wlr_xdg_activation_v1
}

func CreateXDGActivationV1(display Display) XDGActivationV1 {
	p := C.wlr_xdg_activation_v1_create(display.p)
	return XDGActivationV1{p: p}
}

// OnRequestActivate is emitted when a client asks for surface to be
// activated with a token that wlroots knows about. wlroots destroys the
// token once the callback returns.
func (a XDGActivationV1) OnRequestActivate(cb func(XDGActivationV1RequestActivateEvent)) Listener {
	return newListener(&a.p.events.request_activate, func(lis Listener, data unsafe.Pointer) {
		cb(XDGActivationV1RequestActivateEvent{p: (*C.struct_wlr_xdg_activation_v1_request_activate_event)(data)})
	})
}

// OnNewToken is emitted when a client commits a token request that
// wlroots accepted. It isn't emitted for tokens made with CreateToken.
func (a XDGActivationV1) OnNewToken(cb func(XDGActivationTokenV1)) Listener {
	return newListener(&a.p.events.new_token, func(lis Listener, data unsafe.Pointer) {
		cb(XDGActivationTokenV1{p: (*C.struct_wlr_xdg_activation_token_v1)(data)})
	})
}

// CreateToken creates a token that isn't tied to any client, seat or
// surface, such as one to hand to a program that the compositor
// starts.
func (a XDGActivationV1) CreateToken() XDGActivationTokenV1 {
	p := C.wlr_xdg_activation_token_v1_create(a.p)
	return XDGActivationTokenV1{p: p}
}

type XDGActivationV1RequestActivateEvent struct {
	p *C.struct_wlr_xdg_activation_v1_request_activate_event
}

func (e XDGActivationV1RequestActivateEvent) Token() XDGActivationTokenV1 {
	return XDGActivationTokenV1{p: e.p.token}
}

// Surface is the surface that the client wants activated.
func (e XDGActivationV1RequestActivateEvent) Surface() Surface {
	return Surface{p: e.p.surface}
}

type XDGActivationTokenV1 struct {
	p *C.struct_wlr_xdg_activation_token_v1
}

func (t XDGActivationTokenV1) Valid() bool {
	return t.p != nil
}

// Surface is the surface that the client said the token came from. It's
// invalid if the client didn't give one or if the surface is gone.
func (t XDGActivationTokenV1) Surface() Surface {
	return Surface{p: t.p.surface}
}

// Seat is the seat of the input event that the client said the token
// came from. It's invalid if the client didn't give one or if the seat
// is gone.
func (t XDGActivationTokenV1) Seat() Seat {
	return Seat{p: t.p.seat}
}

// Serial is the serial of the input event that the client said the
// token came from. It's only meaningful if Seat is valid.
func (t XDGActivationTokenV1) Serial() uint32 {
	return uint32(t.p.serial)
}

// Name returns the token's string, such as to export in the
// environment of a program that the compositor starts.
func (t XDGActivationTokenV1) Name() string {
	return C.GoString(C.wlr_xdg_activation_token_v1_get_name(t.p))
}

func (t XDGActivationTokenV1) OnDestroy(cb func(XDGActivationTokenV1)) Listener {
	return newListener(&t.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(t)
	})
}
