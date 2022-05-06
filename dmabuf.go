package wlr

/*
#include <wlr/types/wlr_linux_dmabuf_v1.h>
#include <wlr/types/wlr_export_dmabuf_v1.h>
*/
import "C"

import "unsafe"

type LinuxDMABufV1 struct {
	p *C.struct_wlr_linux_dmabuf_v1
}

func CreateLinuxDMABufV1(display Display, renderer Renderer) LinuxDMABufV1 {
	p := C.wlr_linux_dmabuf_v1_create(display.p, renderer.p)
	return LinuxDMABufV1{p: p}
}

func (b LinuxDMABufV1) OnDestroy(cb func(LinuxDMABufV1)) Listener {
	return newListener(&b.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(b)
	})
}

type ExportDMABufManagerV1 struct {
	p *C.struct_wlr_export_dmabuf_manager_v1
}

func CreateExportDMABufV1(display Display) ExportDMABufManagerV1 {
	p := C.wlr_export_dmabuf_manager_v1_create(display.p)
	return ExportDMABufManagerV1{p: p}
}

func (b ExportDMABufManagerV1) OnDestroy(cb func(ExportDMABufManagerV1)) Listener {
	return newListener(&b.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(b)
	})
}
