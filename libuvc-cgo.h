#ifndef LIBUVC_CGO_H
#define LIBUVC_CGO_H

#include <libuvc-binding.h>

typedef void (cgo_uvc_device_callback_t)(uvc_device_t* device, int index, void* ptr);

// device callback go func defined in uvc.go
void go_device_cb(uvc_device_t* device, int index, void* ptr);
void cgo_device_cb(uvc_device_t* device, int index, void *ptr);

// frame callback go func defined in frame.go
void go_frame_cb(uvc_frame_t *frame, void *ptr);
void cgo_frame_cb(uvc_frame_t *frame, void *ptr);

uvc_error_t cgo_uvc_get_device_list(uvc_context_t *ctx, cgo_uvc_device_callback_t *cb, void* ptr);

#endif