package uvc

/*
#include <libuvc-cgo.h>
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/mattn/go-pointer"
)

type UVC struct {
	ctx     *C.uvc_context_t
	devices []*Device
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

//export go_device_cb
func go_device_cb(device *C.uvc_device_t, index C.int, p unsafe.Pointer) {
	uvc := pointer.Restore(p).(*UVC)

	uvc.devices = append(uvc.devices, &Device{dev: device})
}

func (uvc *UVC) GetDevices() ([]*Device, error) {
	uvc.devices = nil

	p := pointer.Save(uvc)
	defer pointer.Unref(p)

	r := C.cgo_uvc_get_device_list(uvc.ctx,
		(*C.cgo_uvc_device_callback_t)(unsafe.Pointer(C.cgo_device_cb)), p)
	if err := newError(ErrorType(r)); err != nil {
		return nil, err
	}

	return uvc.devices, nil
}

// Exit closes the UVC context, shutting down any active devices.
func (uvc *UVC) Exit() {
	C.uvc_exit(uvc.ctx)
}
