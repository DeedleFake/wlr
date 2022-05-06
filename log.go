package wlr

/*
#include <stdlib.h>
#include <stdio.h>

#include <wlr/util/log.h>

void _wlr_log_cb(enum wlr_log_importance importance, char *msg);

static inline void _wlr_log_inner_cb(enum wlr_log_importance importance, const char *fmt, va_list args) {
	char *msg = NULL;
	if (vasprintf(&msg, fmt, args) == -1) {
		return;
	}

	_wlr_log_cb(importance, msg);
	free(msg);
}

static inline void _wlr_log_set_cb(enum wlr_log_importance verbosity, bool is_set) {
	wlr_log_init(verbosity, is_set ? &_wlr_log_inner_cb : NULL);
}
*/
import "C"

type LogImportance uint32

const (
	Silent LogImportance = C.WLR_SILENT
	Error  LogImportance = C.WLR_ERROR
	Info   LogImportance = C.WLR_INFO
	Debug  LogImportance = C.WLR_DEBUG
)

type LogFunc func(importance LogImportance, msg string)

var onLog LogFunc

func InitLog(verbosity LogImportance, cb LogFunc) {
	C._wlr_log_set_cb(C.enum_wlr_log_importance(verbosity), cb != nil)
	onLog = cb
}

//export _wlr_log_cb
func _wlr_log_cb(importance LogImportance, msg *C.char) {
	if onLog != nil {
		onLog(importance, C.GoString(msg))
	}
}
