package wlr

// #include <wlr/types/wlr_layer_shell_v1.h>
import "C"

type LayerShellV1 struct {
	p *C.struct_wlr_layer_shell_v1
}
