package wlr

/*
#include <stdlib.h>
#include <wlr/backend/wayland.h>
#include <wlr/backend/x11.h>
#include <wlr/types/wlr_output.h>
#include <wlr/types/wlr_output_layout.h>
#include <wlr/util/transform.h>
*/
import "C"

import (
	"errors"
	"image"
	"iter"
	"unsafe"
)

type Output struct {
	p *C.struct_wlr_output
}

func (o Output) OnDestroy(cb func(Output)) Listener {
	return newListener(&o.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(o)
	})
}

func (o Output) Name() string {
	return C.GoString(o.p.name)
}

func (o Output) Scale() float32 {
	return float32(o.p.scale)
}

func (o Output) Transform() OutputTransform {
	return OutputTransform(o.p.transform)
}

func (o Output) OnFrame(cb func(Output)) Listener {
	return newListener(&o.p.events.frame, func(lis Listener, data unsafe.Pointer) {
		cb(o)
	})
}

func (o Output) OnRequestState(cb func(Output, OutputState)) Listener {
	return newListener(&o.p.events.request_state, func(lis Listener, data unsafe.Pointer) {
		event := (*C.struct_wlr_output_event_request_state)(data)
		cb(o, OutputState{p: event.state})
	})
}

func (o Output) AddSoftwareCursorsToRenderPass(pass RenderPass, damage image.Rectangle) {
	var cd *C.pixman_region32_t
	if damage != (image.Rectangle{}) {
		t := rectToC(damage)
		cd = &t
	}
	C.wlr_output_add_software_cursors_to_render_pass(o.p, pass.p, cd)
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

func (o Output) InitRender(a Allocator, r Renderer) bool {
	return bool(C.wlr_output_init_render(o.p, a.p, r.p))
}

func (o Output) BeginRenderPass(state OutputState) (RenderPass, error) {
	pass := C.wlr_output_begin_render_pass(o.p, state.p, nil)
	if pass == nil {
		return RenderPass{}, errors.New("can't begin render pass")
	}
	return RenderPass{p: pass}, nil
}

func (o Output) CreateGlobal(display Display) {
	C.wlr_output_create_global(o.p, display.p)
}

func (o Output) DestroyGlobal() {
	C.wlr_output_destroy_global(o.p)
}

func (o Output) TestState(state OutputState) bool {
	return bool(C.wlr_output_test_state(o.p, state.p))
}

func (o Output) CommitState(state OutputState) bool {
	return bool(C.wlr_output_commit_state(o.p, state.p))
}

func (o Output) ScheduleFrame() {
	C.wlr_output_schedule_frame(o.p)
}

func (o Output) Modes() iter.Seq[OutputMode] {
	offset := int(unsafe.Offsetof(C.struct_wlr_output_mode{}.link))
	return func(yield func(OutputMode) bool) {
		seq := listSeq[C.struct_wlr_output_mode](&o.p.modes, offset)
		for mode := range seq {
			if !yield(OutputMode{p: mode}) {
				return
			}
		}
	}
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

func (o Output) PreferredMode() OutputMode {
	p := C.wlr_output_preferred_mode(o.p)
	return OutputMode{p: p}
}

func (o Output) Width() int {
	return int(o.p.width)
}

func (o Output) Height() int {
	return int(o.p.height)
}

func (o Output) Enabled() bool {
	return bool(o.p.enabled)
}

type OutputState struct {
	p *C.struct_wlr_output_state
}

// NewOutputState allocates and initializes an output state on the Go heap.
// Call Finish when done.
func NewOutputState() OutputState {
	state := OutputState{p: &C.struct_wlr_output_state{}}
	C.wlr_output_state_init(state.p)
	return state
}

func (s OutputState) Valid() bool {
	return s.p != nil
}

func (s OutputState) Finish() {
	if s.p == nil {
		return
	}
	C.wlr_output_state_finish(s.p)
}

func (s OutputState) SetEnabled(enabled bool) {
	C.wlr_output_state_set_enabled(s.p, C.bool(enabled))
}

func (s OutputState) SetMode(mode OutputMode) {
	C.wlr_output_state_set_mode(s.p, mode.p)
}

func (s OutputState) SetScale(scale float32) {
	C.wlr_output_state_set_scale(s.p, C.float(scale))
}

func (s OutputState) SetTransform(transform OutputTransform) {
	C.wlr_output_state_set_transform(s.p, C.enum_wl_output_transform(transform))
}

type OutputLayout struct {
	p *C.struct_wlr_output_layout
}

func CreateOutputLayout(display Display) OutputLayout {
	p := C.wlr_output_layout_create(display.p)
	return OutputLayout{p: p}
}

func (l OutputLayout) Destroy() {
	C.wlr_output_layout_destroy(l.p)
}

func (l OutputLayout) AddAuto(output Output) OutputLayoutOutput {
	p := C.wlr_output_layout_add_auto(l.p, output.p)
	return OutputLayoutOutput{p: p}
}

func (l OutputLayout) Add(output Output, lx, ly int) OutputLayoutOutput {
	p := C.wlr_output_layout_add(l.p, output.p, C.int(lx), C.int(ly))
	return OutputLayoutOutput{p: p}
}

func (l OutputLayout) OutputCoords(output Output) (x float64, y float64) {
	var ox, oy C.double
	C.wlr_output_layout_output_coords(l.p, output.p, &ox, &oy)
	return float64(ox), float64(oy)
}

func (l OutputLayout) OutputAt(x, y float64) Output {
	p := C.wlr_output_layout_output_at(l.p, C.double(x), C.double(y))
	return Output{p: p}
}

func (l OutputLayout) Get(output Output) OutputLayoutOutput {
	p := C.wlr_output_layout_get(l.p, output.p)
	return OutputLayoutOutput{p: p}
}

type OutputLayoutOutput struct {
	p *C.struct_wlr_output_layout_output
}

func (o OutputLayoutOutput) Valid() bool {
	return o.p != nil
}

func (o OutputLayoutOutput) X() int {
	return int(o.p.x)
}

func (o OutputLayoutOutput) Y() int {
	return int(o.p.y)
}

type OutputMode struct {
	p *C.struct_wlr_output_mode
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

func (m OutputMode) Valid() bool {
	return m.p != nil
}

type OutputTransform int

const (
	OutputTransformNormal     OutputTransform = C.WL_OUTPUT_TRANSFORM_NORMAL
	OutputTransform90         OutputTransform = C.WL_OUTPUT_TRANSFORM_90
	OutputTransform180        OutputTransform = C.WL_OUTPUT_TRANSFORM_180
	OutputTransform270        OutputTransform = C.WL_OUTPUT_TRANSFORM_270
	OutputTransformFlipped    OutputTransform = C.WL_OUTPUT_TRANSFORM_FLIPPED
	OutputTransformFlipped90  OutputTransform = C.WL_OUTPUT_TRANSFORM_FLIPPED_90
	OutputTransformFlipped180 OutputTransform = C.WL_OUTPUT_TRANSFORM_FLIPPED_180
	OutputTransformFlipped270 OutputTransform = C.WL_OUTPUT_TRANSFORM_FLIPPED_270
)

func (transform OutputTransform) Invert() OutputTransform {
	return OutputTransform(C.wlr_output_transform_invert(C.enum_wl_output_transform(transform)))
}
