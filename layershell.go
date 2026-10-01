package wlr

/*
#include <wlr/types/wlr_layer_shell_v1.h>
*/
import "C"

import "unsafe"

type LayerShellV1 struct {
	p *C.struct_wlr_layer_shell_v1
}

func CreateLayerShellV1(display Display, version uint32) LayerShellV1 {
	p := C.wlr_layer_shell_v1_create(display.p, C.uint32_t(version))
	return LayerShellV1{p: p}
}

func (ls LayerShellV1) OnDestroy(cb func(LayerShellV1)) Listener {
	return newListener(&ls.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(ls)
	})
}

// OnNewSurface is emitted when a client creates a layer surface. Its
// output may be invalid, in which case the callback has to assign one
// with SetOutput before it returns.
func (ls LayerShellV1) OnNewSurface(cb func(LayerSurfaceV1)) Listener {
	return newListener(&ls.p.events.new_surface, func(lis Listener, data unsafe.Pointer) {
		cb(LayerSurfaceV1{p: (*C.struct_wlr_layer_surface_v1)(data)})
	})
}

type LayerSurfaceV1 struct {
	p *C.struct_wlr_layer_surface_v1
}

func (s LayerSurfaceV1) Surface() Surface {
	return Surface{
		p: s.p.surface,
	}
}

func (s LayerSurfaceV1) Output() Output {
	return Output{p: s.p.output}
}

func (s LayerSurfaceV1) SetOutput(output Output) {
	s.p.output = output.p
}

func (s LayerSurfaceV1) Namespace() string {
	return C.GoString(s.p.namespace)
}

// Initialized returns true once the client has made its initial
// commit. Configure can't be called before that.
func (s LayerSurfaceV1) Initialized() bool {
	return bool(s.p.initialized)
}

// InitialCommit returns true if the commit currently being handled is
// the surface's initial commit.
func (s LayerSurfaceV1) InitialCommit() bool {
	return bool(s.p.initial_commit)
}

func (s LayerSurfaceV1) Current() LayerSurfaceV1State {
	return LayerSurfaceV1State{v: s.p.current}
}

func (s LayerSurfaceV1) Pending() LayerSurfaceV1State {
	return LayerSurfaceV1State{v: s.p.pending}
}

// Configure sends a configure with the given size and returns its
// serial. The surface must be initialized.
func (s LayerSurfaceV1) Configure(width, height uint32) uint32 {
	return uint32(C.wlr_layer_surface_v1_configure(s.p, C.uint32_t(width), C.uint32_t(height)))
}

// Destroy sends the closed event to the client and destroys the layer
// surface.
func (s LayerSurfaceV1) Destroy() {
	C.wlr_layer_surface_v1_destroy(s.p)
}

// ExclusiveEdge returns the edge that the exclusive zone applies to,
// or EdgeNone if there isn't one.
func (s LayerSurfaceV1) ExclusiveEdge() Edges {
	return Edges(C.wlr_layer_surface_v1_get_exclusive_edge(s.p))
}

func (s LayerSurfaceV1) OnDestroy(cb func(LayerSurfaceV1)) Listener {
	return newListener(&s.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(s)
	})
}

type LayerSurfaceV1State struct {
	v C.struct_wlr_layer_surface_v1_state
}

// Committed returns the fields that the last commit changed.
func (s LayerSurfaceV1State) Committed() LayerSurfaceV1StateField {
	return LayerSurfaceV1StateField(s.v.committed)
}

func (s LayerSurfaceV1State) Anchor() LayerSurfaceV1Anchor {
	return LayerSurfaceV1Anchor(s.v.anchor)
}

func (s LayerSurfaceV1State) ExclusiveZone() int32 {
	return int32(s.v.exclusive_zone)
}

func (s LayerSurfaceV1State) Margin() (top, right, bottom, left int32) {
	m := s.v.margin
	return int32(m.top), int32(m.right), int32(m.bottom), int32(m.left)
}

func (s LayerSurfaceV1State) KeyboardInteractive() LayerSurfaceV1KeyboardInteractivity {
	return LayerSurfaceV1KeyboardInteractivity(s.v.keyboard_interactive)
}

func (s LayerSurfaceV1State) DesiredWidth() uint32 {
	return uint32(s.v.desired_width)
}

func (s LayerSurfaceV1State) DesiredHeight() uint32 {
	return uint32(s.v.desired_height)
}

func (s LayerSurfaceV1State) Layer() LayerShellV1Layer {
	return LayerShellV1Layer(s.v.layer)
}

type LayerShellV1Layer int

const (
	LayerShellV1LayerBackground LayerShellV1Layer = C.ZWLR_LAYER_SHELL_V1_LAYER_BACKGROUND
	LayerShellV1LayerBottom     LayerShellV1Layer = C.ZWLR_LAYER_SHELL_V1_LAYER_BOTTOM
	LayerShellV1LayerTop        LayerShellV1Layer = C.ZWLR_LAYER_SHELL_V1_LAYER_TOP
	LayerShellV1LayerOverlay    LayerShellV1Layer = C.ZWLR_LAYER_SHELL_V1_LAYER_OVERLAY
)

type LayerSurfaceV1Anchor uint32

const (
	LayerSurfaceV1AnchorTop    LayerSurfaceV1Anchor = C.ZWLR_LAYER_SURFACE_V1_ANCHOR_TOP
	LayerSurfaceV1AnchorBottom LayerSurfaceV1Anchor = C.ZWLR_LAYER_SURFACE_V1_ANCHOR_BOTTOM
	LayerSurfaceV1AnchorLeft   LayerSurfaceV1Anchor = C.ZWLR_LAYER_SURFACE_V1_ANCHOR_LEFT
	LayerSurfaceV1AnchorRight  LayerSurfaceV1Anchor = C.ZWLR_LAYER_SURFACE_V1_ANCHOR_RIGHT
)

type LayerSurfaceV1KeyboardInteractivity uint32

const (
	LayerSurfaceV1KeyboardInteractivityNone      LayerSurfaceV1KeyboardInteractivity = C.ZWLR_LAYER_SURFACE_V1_KEYBOARD_INTERACTIVITY_NONE
	LayerSurfaceV1KeyboardInteractivityExclusive LayerSurfaceV1KeyboardInteractivity = C.ZWLR_LAYER_SURFACE_V1_KEYBOARD_INTERACTIVITY_EXCLUSIVE
	LayerSurfaceV1KeyboardInteractivityOnDemand  LayerSurfaceV1KeyboardInteractivity = C.ZWLR_LAYER_SURFACE_V1_KEYBOARD_INTERACTIVITY_ON_DEMAND
)

type LayerSurfaceV1StateField uint32

const (
	LayerSurfaceV1StateDesiredSize           LayerSurfaceV1StateField = C.WLR_LAYER_SURFACE_V1_STATE_DESIRED_SIZE
	LayerSurfaceV1StateAnchor                LayerSurfaceV1StateField = C.WLR_LAYER_SURFACE_V1_STATE_ANCHOR
	LayerSurfaceV1StateExclusiveZone         LayerSurfaceV1StateField = C.WLR_LAYER_SURFACE_V1_STATE_EXCLUSIVE_ZONE
	LayerSurfaceV1StateMargin                LayerSurfaceV1StateField = C.WLR_LAYER_SURFACE_V1_STATE_MARGIN
	LayerSurfaceV1StateKeyboardInteractivity LayerSurfaceV1StateField = C.WLR_LAYER_SURFACE_V1_STATE_KEYBOARD_INTERACTIVITY
	LayerSurfaceV1StateLayer                 LayerSurfaceV1StateField = C.WLR_LAYER_SURFACE_V1_STATE_LAYER
	LayerSurfaceV1StateExclusiveEdge         LayerSurfaceV1StateField = C.WLR_LAYER_SURFACE_V1_STATE_EXCLUSIVE_EDGE
)
