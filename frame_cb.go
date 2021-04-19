package uvc

/*
#include <libuvc-binding.h>

void goFrameCB(uvc_frame_t *frame, void *ptr);

// The gateway function
void frame_cb(uvc_frame_t *frame, void *ptr) {
	goFrameCB(frame, ptr);
}
*/
import "C"
