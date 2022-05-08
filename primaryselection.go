package wlr

/*
#include <wlr/types/wlr_primary_selection_v1.h>
*/
import "C"
import "unsafe"

type PrimarySelectionV1DeviceManager struct {
	p *C.struct_wlr_primary_selection_v1_device_manager
}

func CreatePrimarySelectionV1DeviceManager(display Display) PrimarySelectionV1DeviceManager {
	p := C.wlr_primary_selection_v1_device_manager_create(display.p)
	return PrimarySelectionV1DeviceManager{p: p}
}

func (m PrimarySelectionV1DeviceManager) OnDestroy(cb func(PrimarySelectionV1DeviceManager)) Listener {
	return newListener(&m.p.events.destroy, func(lis Listener, data unsafe.Pointer) {
		cb(m)
	})
}
