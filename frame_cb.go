package uvc

/*
#include <libuvc/libuvc.h>

void goFrameCB(uvc_frame_t *frame, void *ptr);

// The gateway function
void frame_cb(uvc_frame_t *frame, void *ptr) {
	goFrameCB(frame, ptr);
}
*/
import "C"
