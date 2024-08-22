package wlr

/*
#include <stdlib.h>

#include <wlr/types/wlr_cursor.h>
#include <wlr/types/wlr_xcursor_manager.h>
*/
import "C"

import (
	"iter"
	"unsafe"
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

func CreateXCursorManager(name string, size uint32) XCursorManager {
	var cname *C.char
	if name != "" {
		cname = C.CString(name)
	}

	p := C.wlr_xcursor_manager_create(cname, C.uint32_t(size))
	return XCursorManager{p: p}
}

func (m XCursorManager) Destroy() {
	if m.p.name != nil {
		C.free(unsafe.Pointer(m.p.name))
	}

	C.wlr_xcursor_manager_destroy(m.p)
}

func (m XCursorManager) Load(scale float64) {
	C.wlr_xcursor_manager_load(m.p, C.float(scale))
}

func (m XCursorManager) GetXCursor(name string, scale float32) XCursor {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))

	p := C.wlr_xcursor_manager_get_xcursor(m.p, s, C.float(scale))
	return XCursor{p: p}
}

func (c XCursor) Image(i int) XCursorImage {
	n := c.ImageCount()
	slice := (*[1 << 30]*C.struct_wlr_xcursor_image)(unsafe.Pointer(c.p.images))[:n:n]
	return XCursorImage{p: slice[i]}
}

func (c XCursor) Images() iter.Seq[XCursorImage] {
	return func(yield func(XCursorImage) bool) {
		count := c.ImageCount()
		for i := 0; i < count; i++ {
			if !yield(c.Image(i)) {
				return
			}
		}
	}
}

func (c XCursor) ImageCount() int {
	return int(c.p.image_count)
}

func (c XCursor) Name() string {
	return C.GoString(c.p.name)
}
