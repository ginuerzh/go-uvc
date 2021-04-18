package uvc

/*
#include <libuvc/libuvc.h>
#include <stdio.h>

static inline const uvc_frame_desc_t* uvc_get_frame_desc(uvc_format_desc_t* format_desc) {
	return format_desc->frame_descs;
}
*/
import "C"

import (
	"errors"
	"unsafe"
)

// Auto-exposure mode.
type AEMode uint8

const (
	// manual exposure time, manual iris
	AEModeManual AEMode = 1
	// auto exposure time, auto iris
	AEModeAuto AEMode = 2
	// manual exposure time, auto iris
	AEModeShutterPriority AEMode = 4
	// auto exposure time, manual iris
	AEModeAperturePriority AEMode = 8
)

// VideoStreaming interface descriptor subtype.
type VSDescSubType C.enum_uvc_vs_desc_subtype

const (
	VS_UNDEFINED           VSDescSubType = C.UVC_VS_UNDEFINED
	VS_INPUT_HEADER        VSDescSubType = C.UVC_VS_INPUT_HEADER
	VS_OUTPUT_HEADER       VSDescSubType = C.UVC_VS_OUTPUT_HEADER
	VS_STILL_IMAGE_FRAME   VSDescSubType = C.UVC_VS_STILL_IMAGE_FRAME
	VS_FORMAT_UNCOMPRESSED VSDescSubType = C.UVC_VS_FORMAT_UNCOMPRESSED
	VS_FRAME_UNCOMPRESSED  VSDescSubType = C.UVC_VS_FRAME_UNCOMPRESSED
	VS_FORMAT_MJPEG        VSDescSubType = C.UVC_VS_FORMAT_MJPEG
	VS_FRAME_MJPEG         VSDescSubType = C.UVC_VS_FRAME_MJPEG
	VS_FORMAT_MPEG2TS      VSDescSubType = C.UVC_VS_FORMAT_MPEG2TS
	VS_FORMAT_DV           VSDescSubType = C.UVC_VS_FORMAT_DV
	VS_COLORFORMAT         VSDescSubType = C.UVC_VS_COLORFORMAT
	VS_FORMAT_FRAME_BASED  VSDescSubType = C.UVC_VS_FORMAT_FRAME_BASED
	VS_FRAME_FRAME_BASED   VSDescSubType = C.UVC_VS_FRAME_FRAME_BASED
	VS_FORMAT_STREAM_BASED VSDescSubType = C.UVC_VS_FORMAT_STREAM_BASED
)

type Device struct {
	dev    *C.uvc_device_t
	handle *C.uvc_device_handle_t
}

// Ope opens a UVC device.
func (dev *Device) Open() error {
	r := C.uvc_open(dev.dev, &dev.handle)
	return newError(ErrorType(r))
}

func (dev *Device) SetAEMode(mode AEMode) error {
	r := C.uvc_set_ae_mode(dev.handle, C.uchar(mode))
	return newError(ErrorType(r))
}

// GetBusNumber gets the number of the bus to which the device is attached.
func (dev *Device) GetBusNumber() uint8 {
	return uint8(C.uvc_get_bus_number(dev.dev))
}

// GetAddress gets the number assigned to the device within its bus.
func (dev *Device) GetAddress() uint8 {
	return uint8(C.uvc_get_device_address(dev.dev))
}

func (dev *Device) GetFormatDesc() *FormatDescriptor {
	desc := C.uvc_get_format_descs(dev.handle)
	if desc == nil {
		return nil
	}
	return &FormatDescriptor{
		desc:                desc,
		Subtype:             VSDescSubType(desc.bDescriptorSubtype),
		FormatIndex:         uint8(desc.bFormatIndex),
		NumFrameDescriptors: uint8(desc.bNumFrameDescriptors),
		BitsPerPixel:        *(*uint8)(unsafe.Pointer(&desc.anon1[0])),
		Flags:               *(*uint8)(unsafe.Pointer(&desc.anon1[0])),
		DefaultFrameIndex:   uint8(desc.bDefaultFrameIndex),
		AspectRatioX:        uint8(desc.bAspectRatioX),
		AspectRatioY:        uint8(desc.bAspectRatioY),
		InterlaceFlags:      uint8(desc.bmInterlaceFlags),
		CopyProtect:         uint8(desc.bCopyProtect),
		VariableSize:        uint8(desc.bVariableSize),
	}
}

func (dev *Device) Diag() {
	C.uvc_print_diag(dev.handle, C.stderr)
}

func (dev *Device) Info() (*DeviceInfo, error) {
	if dev.handle == nil {
		return nil, errors.New("device closed")
	}

	var desc *C.uvc_device_descriptor_t
	r := C.uvc_get_device_descriptor(dev.dev, &desc)
	if err := newError(ErrorType(r)); err != nil {
		return nil, err
	}
	defer C.uvc_free_device_descriptor(desc)

	info := &DeviceInfo{
		VendorID:     uint16(desc.idVendor),
		ProductID:    uint16(desc.idProduct),
		BcdUVC:       uint16(desc.bcdUVC),
		SerialNumber: C.GoString(desc.serialNumber),
		Manufacturer: C.GoString(desc.manufacturer),
		Product:      C.GoString(desc.product),
	}
	return info, nil
}

// GetStream gets a negotiated streaming control block for some common parameters.
func (dev *Device) GetStream(format FrameFormat, width, height, fps int) (*Stream, error) {
	var ctrl C.uvc_stream_ctrl_t
	r := C.uvc_get_stream_ctrl_format_size(dev.handle, &ctrl,
		C.enum_uvc_frame_format(format),
		C.int(width), C.int(height), C.int(fps))
	if err := newError(ErrorType(r)); err != nil {
		return nil, err
	}
	return &Stream{
		devh: dev.handle,
		ctrl: ctrl,
		fc:   make(chan *Frame),
	}, nil
}

// Ref increments the reference count for a device.
func (dev *Device) Ref() {
	C.uvc_ref_device(dev.dev)
}

// Unref decrements the reference count for a device.
func (dev *Device) Unref() {
	C.uvc_unref_device(dev.dev)
}

// Close closes a device.
// Ends any stream that's in progress.
// The device handle and frame structures will be invalidated.
func (dev *Device) Close() error {
	C.uvc_close(dev.handle)
	return nil
}

type DeviceInfo struct {
	VendorID     uint16
	ProductID    uint16
	BcdUVC       uint16
	SerialNumber string
	Manufacturer string
	Product      string
}

// device format descriptor.
// A "format" determines a stream's image type (e.g., raw YUYV or JPEG),
// and includes many "frame" configurations.
type FormatDescriptor struct {
	desc *C.uvc_format_desc_t
	// Type of image stream, such as JPEG or uncompressed.
	Subtype VSDescSubType
	// Identifier of this format within the VS interface's format list
	FormatIndex         uint8
	NumFrameDescriptors uint8
	// BPP for uncompressed stream
	BitsPerPixel uint8
	// Flags for JPEG stream
	Flags             uint8
	DefaultFrameIndex uint8
	AspectRatioX      uint8
	AspectRatioY      uint8
	InterlaceFlags    uint8
	CopyProtect       uint8
	VariableSize      uint8
}

func (d *FormatDescriptor) FrameDesc() *FrameDescriptor {
	desc := C.uvc_get_frame_desc(d.desc)
	if desc == nil {
		return nil
	}
	return &FrameDescriptor{
		desc:                    desc,
		Subtype:                 VSDescSubType(desc.bDescriptorSubtype),
		FrameIndex:              uint8(desc.bFrameIndex),
		Capabilities:            uint8(desc.bmCapabilities),
		Width:                   uint16(desc.wWidth),
		Height:                  uint16(desc.wHeight),
		MinBitRate:              uint32(desc.dwMinBitRate),
		MaxBitRate:              uint32(desc.dwMaxBitRate),
		MaxVideoFrameBufferSize: uint32(desc.dwMaxVideoFrameBufferSize),
		DefaultFrameInterval:    uint32(desc.dwDefaultFrameInterval),
		MinFrameInterval:        uint32(desc.dwMinFrameInterval),
		MaxFrameInterval:        uint32(desc.dwMaxFrameInterval),
		FrameIntervalStep:       uint32(desc.dwFrameIntervalStep),
		FrameIntervalType:       uint8(desc.bFrameIntervalType),
		BytesPerLine:            uint32(desc.dwBytesPerLine),
		// Intervals: []uint32,
	}
}

// Frame descriptor.
// A "frame" is a configuration of a streaming format
// for a particular image size at one of possibly several available frame rates.
type FrameDescriptor struct {
	desc *C.uvc_frame_desc_t
	// Type of frame, such as JPEG frame or uncompressed frame
	Subtype VSDescSubType
	// Index of the frame within the list of specs available for this format
	FrameIndex   uint8
	Capabilities uint8
	// Image width
	Width uint16
	// Image height
	Height uint16
	// Bitrate of corresponding stream at minimal frame rate
	MinBitRate uint32
	// Bitrate of corresponding stream at maximal frame rate
	MaxBitRate uint32
	// Maximum number of bytes for a video frame
	MaxVideoFrameBufferSize uint32
	// Default frame interval (in 100ns units)
	DefaultFrameInterval uint32
	// Minimum frame interval for continuous mode (100ns units)
	MinFrameInterval uint32
	// Maximum frame interval for continuous mode (100ns units)
	MaxFrameInterval uint32
	// Granularity of frame interval range for continuous mode (100ns)
	FrameIntervalStep uint32
	// Frame intervals
	FrameIntervalType uint8
	// Number of bytes per line
	BytesPerLine uint32
	// Available frame rates, zero-terminated (in 100ns units)
	Intervals []uint32
}
