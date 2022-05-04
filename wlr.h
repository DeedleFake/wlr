#ifndef _GO_WLROOTS_H
#define _GO_WLROOTS_H

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

void _wlr_xdg_surface_for_each_cb(struct wlr_surface *surface, int sx, int sy, void *data);

static inline void _wlr_xdg_surface_for_each_surface(struct wlr_xdg_surface *surface, void *user_data) {
 	wlr_xdg_surface_for_each_surface(surface, &_wlr_xdg_surface_for_each_cb, user_data);
}

void _wl_listener_cb(struct wl_listener *listener, void *data);

static inline void _wl_listener_set_cb(struct wl_listener *listener) {
 	listener->notify = &_wl_listener_cb;
}

#endif
