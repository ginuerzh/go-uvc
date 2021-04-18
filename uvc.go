package uvc

// #cgo LDFLAGS: -luvc
// #include <libuvc/libuvc.h>
// #include <stdlib.h>
import "C"
import "unsafe"

type UVC struct {
	ctx *C.uvc_context_t
}

// Init initializes a UVC service context.
// Libuvc will set up its own libusb context.
func (uvc *UVC) Init() error {
	res := C.uvc_init(&uvc.ctx, nil)
	return newError(ErrorType(res))
}

// FindDevice finds a device identified by vendor vid, product pid and/or serial number sn.
func (uvc *UVC) FindDevice(vid, pid int, sn string) (*Device, error) {
	var csn *C.char
	if sn != "" {
		csn = C.CString(sn)
		defer C.free(unsafe.Pointer(csn))
	}

	var dev *C.uvc_device_t
	res := C.uvc_find_device(uvc.ctx, &dev, C.int(vid), C.int(pid), csn)
	if err := newError(ErrorType(res)); err != nil {
		return nil, err
	}
	return &Device{dev: dev}, nil
}

// Exit closes the UVC context, shutting down any active devices.
func (uvc *UVC) Exit() {
	C.uvc_exit(uvc.ctx)
}
