package xkb

//go:generate go run ./internal/symsgen syms.go

/*
#cgo pkg-config: xkbcommon

#include <stdlib.h>
#include <xkbcommon/xkbcommon.h>
*/
import "C"

import "unsafe"

type KeyCode uint32

type KeySymFlags uint32

const (
	KeySymNoFlags         KeySymFlags = C.XKB_KEYSYM_NO_FLAGS
	KeySymCaseInsensitive KeySymFlags = C.XKB_KEYSYM_CASE_INSENSITIVE
)

func SymFromName(name string, flags KeySymFlags) KeySym {
	s := C.CString(name)
	sym := C.xkb_keysym_from_name(s, C.enum_xkb_keysym_flags(flags))
	C.free(unsafe.Pointer(s))
	return KeySym(sym)
}

type ContextFlags int

const (
	ContextNoFlags            ContextFlags = C.XKB_CONTEXT_NO_FLAGS
	ContextNoDefaultincludes               = C.XKB_CONTEXT_NO_DEFAULT_INCLUDES
	ContextNoEnvironmentNames              = C.XKB_CONTEXT_NO_ENVIRONMENT_NAMES
)

type Context struct {
	p *C.struct_xkb_context
}

func NewContext(flags ContextFlags) Context {
	p := C.xkb_context_new(C.XKB_CONTEXT_NO_FLAGS)
	return Context{p: p}
}

func (c Context) Unref() {
	C.xkb_context_unref(c.p)
}

type State struct {
	p *C.struct_xkb_state
}

func (s State) Syms(keyCode KeyCode) []KeySym {
	var syms *C.xkb_keysym_t
	n := int(C.xkb_state_key_get_syms(s.p, C.uint32_t(keyCode), &syms))
	if n == 0 || syms == nil {
		return nil
	}
	slice := unsafe.Slice(syms, n)

	res := make([]KeySym, n)
	for i := range n {
		res[i] = KeySym(slice[i])
	}

	return res
}

type RuleNames struct {
	Rules   string
	Model   string
	Layout  string
	Variant string
	Options string
}

func (rn *RuleNames) toC() *C.struct_xkb_rule_names {
	if rn == nil {
		return nil
	}

	return &C.struct_xkb_rule_names{
		rules:   C.CString(rn.Rules),
		model:   C.CString(rn.Model),
		layout:  C.CString(rn.Layout),
		variant: C.CString(rn.Variant),
		options: C.CString(rn.Options),
	}
}
