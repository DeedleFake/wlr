package xkb

/*
#include <xkbcommon/xkbcommon.h>
*/
import "C"

type Keymap struct {
	p *C.struct_xkb_keymap
}

func NewKeymapFromNames(ctx Context, rules *RuleNames, flags KeymapCompileFlags) Keymap {
	p := C.xkb_keymap_new_from_names(ctx.p, rules.toC(), C.enum_xkb_keymap_compile_flags(flags))
	return Keymap{p: p}
}

func (m Keymap) Unref() {
	C.xkb_keymap_unref(m.p)
}

type KeymapCompileFlags int

const (
	KeymapCompileNoFlags KeymapCompileFlags = C.XKB_KEYMAP_COMPILE_NO_FLAGS
)
