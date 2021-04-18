package uvc

/*
#include <libuvc/libuvc.h>

void frame_cb(uvc_frame_t *frame, void *ptr);
*/
import "C"
import (
	"bytes"
	"fmt"
	"unsafe"

	"github.com/mattn/go-pointer"
)

type Stream struct {
	devh   *C.uvc_device_handle_t
	ctrl   C.uvc_stream_ctrl_t
	handle *C.uvc_stream_handle_t
	fc     chan *Frame
	p      unsafe.Pointer
}

// Open opens a new video stream.
func (s *Stream) Open() error {
	r := C.uvc_stream_open_ctrl(s.devh, &s.handle, &s.ctrl)
	return newError(ErrorType(r))
}

// Start begin streaming video from the device into the callback function.
func (s *Stream) Start() (<-chan *Frame, error) {
	s.p = pointer.Save(s.fc)
	r := C.uvc_stream_start(s.handle,
		(*C.uvc_frame_callback_t)(unsafe.Pointer(C.frame_cb)), s.p, 0)
	if err := newError(ErrorType(r)); err != nil {
		return nil, err
	}

	return s.fc, nil
}

func (s *Stream) Stop() error {
	pointer.Unref(s.p)
	r := C.uvc_stream_stop(s.handle)
	return newError(ErrorType(r))
}

func (s *Stream) Close() error {
	C.uvc_stream_close(s.handle)
	return nil
}

//export goFrameCB
func goFrameCB(frame *C.struct_uvc_frame, p unsafe.Pointer) {
	fc := pointer.Restore(p).(chan *Frame)
	fr := &Frame{
		Width:  int(frame.width),
		Height: int(frame.height),
	}
	select {
	case fc <- fr:
	default:
	}
}

func (s *Stream) Ctrl() *StreamCtrl {
	return &StreamCtrl{
		Hint:                   uint16(s.ctrl.bmHint),
		FormatIndex:            uint8(s.ctrl.bFormatIndex),
		FrameIndex:             uint8(s.ctrl.bFrameIndex),
		FrameInterval:          uint32(s.ctrl.dwFrameInterval),
		KeyFrameRate:           uint16(s.ctrl.wKeyFrameRate),
		PFrameRate:             uint16(s.ctrl.wPFrameRate),
		CompQuality:            uint16(s.ctrl.wCompQuality),
		CompWindowSize:         uint16(s.ctrl.wCompWindowSize),
		Delay:                  uint16(s.ctrl.wDelay),
		MaxVideoFrameSize:      uint32(s.ctrl.dwMaxVideoFrameSize),
		MaxPayloadTransferSize: uint32(s.ctrl.dwMaxPayloadTransferSize),
		ClockFrequency:         uint32(s.ctrl.dwClockFrequency),
		FramingInfo:            uint8(s.ctrl.bmFramingInfo),
		PreferredVersion:       uint8(s.ctrl.bPreferredVersion),
		MinVersion:             uint8(s.ctrl.bMinVersion),
		MaxVersion:             uint8(s.ctrl.bMaxVersion),
		InterfaceNumber:        uint8(s.ctrl.bInterfaceNumber),
	}
}

type StreamCtrl struct {
	Hint                   uint16
	FormatIndex            uint8
	FrameIndex             uint8
	FrameInterval          uint32
	KeyFrameRate           uint16
	PFrameRate             uint16
	CompQuality            uint16
	CompWindowSize         uint16
	Delay                  uint16
	MaxVideoFrameSize      uint32
	MaxPayloadTransferSize uint32
	ClockFrequency         uint32
	FramingInfo            uint8
	PreferredVersion       uint8
	MinVersion             uint8
	MaxVersion             uint8
	InterfaceNumber        uint8
}

func (sc *StreamCtrl) String() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "\nHint: %04x\n", sc.Hint)
	fmt.Fprintf(buf, "FormatIndex: %d\n", sc.FormatIndex)
	fmt.Fprintf(buf, "FrameIndex: %d\n", sc.FrameIndex)
	fmt.Fprintf(buf, "FrameInterval: %d\n", sc.FrameInterval)
	fmt.Fprintf(buf, "KeyFrameRate: %d\n", sc.KeyFrameRate)
	fmt.Fprintf(buf, "PFrameRate: %d\n", sc.PFrameRate)
	fmt.Fprintf(buf, "CompQuality: %d\n", sc.CompQuality)
	fmt.Fprintf(buf, "CompWindowSize: %d\n", sc.CompWindowSize)
	fmt.Fprintf(buf, "Delay: %d\n", sc.Delay)
	fmt.Fprintf(buf, "MaxVideoFrameSize: %d\n", sc.MaxVideoFrameSize)
	fmt.Fprintf(buf, "MaxPayloadTransferSize: %d\n", sc.MaxPayloadTransferSize)
	fmt.Fprintf(buf, "ClockFrequency: %d\n", sc.ClockFrequency)
	fmt.Fprintf(buf, "FramingInfo: %d\n", sc.FramingInfo)
	fmt.Fprintf(buf, "PreferredVersion: %d\n", sc.PreferredVersion)
	fmt.Fprintf(buf, "MinVersion: %d\n", sc.MinVersion)
	fmt.Fprintf(buf, "MaxVersion: %d\n", sc.MaxVersion)
	fmt.Fprintf(buf, "InterfaceNumber: %d\n", sc.InterfaceNumber)
	return buf.String()
}
