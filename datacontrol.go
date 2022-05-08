package wlr

/*
#include <wlr/types/wlr_data_control_v1.h>
*/
import "C"

import "unsafe"

type DataControlManagerV1 struct {
	p *C.struct_wlr_data_control_manager_v1
}

func CreateDataControlManagerV1(display Display) DataControlManagerV1 {
	p := C.wlr_data_control_manager_v1_create(display.p)
	return DataControlManagerV1{p: p}
}

func (m DataControlManagerV1) OnDestroy(cb func(DataControlManagerV1)) Listener {
	return newListener(&m.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(m)
	})
}

func (m DataControlManagerV1) OnNewDevice(cb func(DataControlDeviceV1)) Listener {
	return newListener(&m.p.events.new_device, func(lis Listener, data unsafe.Pointer) {
		cb(DataControlDeviceV1{p: (*C.struct_wlr_data_control_device_v1)(data)})
	})
}

type DataControlDeviceV1 struct {
	p *C.struct_wlr_data_control_device_v1
}

func (d DataControlDeviceV1) Destroy() {
	C.wlr_data_control_device_v1_destroy(d.p)
}
