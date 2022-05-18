package wlr

/*
#include <wlr/render/allocator.h>
*/
import "C"

type Allocator struct {
	p *C.struct_wlr_allocator
}

func AutocreateAllocator(backend Backend, renderer Renderer) Allocator {
	p := C.wlr_allocator_autocreate(backend.p, renderer.p)
	return Allocator{p: p}
}

func (a Allocator) Valid() bool {
	return a.p != nil
}

func (a Allocator) Destroy() {
	C.wlr_allocator_destroy(a.p)
}

//func (a Allocator) OnDestroy(cb func(Allocator)) Listener {
//	return newListener(&a.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
//		cb(a)
//	})
//}

//func (a Allocator)CreateBuffer(width, height int, format *DRMFormat) Buffer {
//	p := C.wlr_allocator_create_buffer(a.p, C.int(width), C.int(height), format.toC())
//	return Buffer{p: p}
//}
