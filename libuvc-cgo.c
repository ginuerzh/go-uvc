#include <libuvc-cgo.h>

// The callback gateway function for device callback.
void cgo_device_cb(uvc_device_t* device, int index, void* ptr) {
	go_device_cb(device, index, ptr);
}

// The callback gateway function for frame callback.
void cgo_frame_cb(uvc_frame_t *frame, void *ptr) {
	go_frame_cb(frame, ptr);
}

// Just like uvc_get_device_list
uvc_error_t cgo_uvc_get_device_list(uvc_context_t* ctx, cgo_uvc_device_callback_t* cb, void* ptr) {
	uvc_error_t ret;
	uvc_device_t **devices;
	uvc_device_t *device;
	int dev_idx;

	ret = uvc_get_device_list(ctx, &devices);
	if (ret != UVC_SUCCESS) {
		UVC_EXIT(ret);
		return ret;
	}

	dev_idx = 0;
	while((device = devices[dev_idx]) != NULL) {
		cb(device, dev_idx, ptr);
		dev_idx++;
	}

	uvc_free_device_list(devices, 1);

	UVC_EXIT(UVC_SUCCESS);
	return UVC_SUCCESS;
}
