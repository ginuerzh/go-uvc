package uvc

/*
#include <libuvc-cgo.h>
*/
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"unsafe"

	"github.com/mattn/go-pointer"
)

var (
	ErrStreamClosed = errors.New("stream closed")
)

type Stream struct {
	devh   *C.uvc_device_handle_t
	handle *C.uvc_stream_handle_t
	ctrl   C.uvc_stream_ctrl_t
	fc     chan *Frame
	p      unsafe.Pointer
	mu     sync.RWMutex
}

// Open opens a new video stream.
func (s *Stream) Open() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.handle != nil {
		return nil
	}

	r := C.uvc_stream_open_ctrl(s.devh, &s.handle, &s.ctrl)
	if err := newError(ErrorType(r)); err != nil {
		return err
	}

	s.fc = make(chan *Frame)
	return nil
}

// Start begins streaming video from the device into frame channel.
func (s *Stream) Start() (<-chan *Frame, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.handle == nil {
		return nil, ErrStreamClosed
	}

	s.p = pointer.Save(s.fc)
	r := C.uvc_stream_start(s.handle,
		(*C.uvc_frame_callback_t)(unsafe.Pointer(C.cgo_frame_cb)), s.p, 0)
	if err := newError(ErrorType(r)); err != nil {
		return nil, err
	}

	return s.fc, nil
}

func (s *Stream) Stop() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.handle == nil {
		return ErrStreamClosed
	}

	pointer.Unref(s.p)
	r := C.uvc_stream_stop(s.handle)
	return newError(ErrorType(r))
}

func (s *Stream) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.handle == nil {
		return nil
	}

	C.uvc_stream_close(s.handle)
	s.handle = nil
	close(s.fc)

	return nil
}

func (s *Stream) IsClosed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.handle == nil
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
	fmt.Fprintf(buf, "Hint: %04x\n", sc.Hint)
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
