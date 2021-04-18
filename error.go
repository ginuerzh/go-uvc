package uvc

import (
	"errors"
	"fmt"
)

// #include <libuvc/libuvc.h>
import "C"

type ErrorType C.enum_uvc_error

const (
	SUCCESS               ErrorType = C.UVC_SUCCESS
	ERROR_IO              ErrorType = C.UVC_ERROR_IO
	ERROR_INVALID_PARAM   ErrorType = C.UVC_ERROR_INVALID_PARAM
	ERROR_ACCESS          ErrorType = C.UVC_ERROR_ACCESS
	ERROR_NO_DEVICE       ErrorType = C.UVC_ERROR_NO_DEVICE
	ERROR_NOT_FOUND       ErrorType = C.UVC_ERROR_NOT_FOUND
	ERROR_BUSY            ErrorType = C.UVC_ERROR_BUSY
	ERROR_TIMEOUT         ErrorType = C.UVC_ERROR_TIMEOUT
	ERROR_OVERFLOW        ErrorType = C.UVC_ERROR_OVERFLOW
	ERROR_PIPE            ErrorType = C.UVC_ERROR_PIPE
	ERROR_INTERRUPTED     ErrorType = C.UVC_ERROR_INTERRUPTED
	ERROR_NO_MEM          ErrorType = C.UVC_ERROR_NO_MEM
	ERROR_NOT_SUPPORTED   ErrorType = C.UVC_ERROR_NOT_SUPPORTED
	ERROR_INVALID_DEVICE  ErrorType = C.UVC_ERROR_INVALID_DEVICE
	ERROR_INVALID_MODE    ErrorType = C.UVC_ERROR_INVALID_MODE
	ERROR_CALLBACK_EXISTS ErrorType = C.UVC_ERROR_CALLBACK_EXISTS
	ERROR_OTHER           ErrorType = C.UVC_ERROR_OTHER
)

func newError(et ErrorType) error {
	switch et {
	case SUCCESS:
		return nil
	case ERROR_IO:
		return errors.New("io error")
	case ERROR_INVALID_PARAM:
		return errors.New("invalid parameter")
	case ERROR_ACCESS:
		return errors.New("access denied")
	case ERROR_NO_DEVICE:
		return errors.New("no such device")
	case ERROR_NOT_FOUND:
		return errors.New("entity not found")
	case ERROR_BUSY:
		return errors.New("resource busy")
	case ERROR_TIMEOUT:
		return errors.New("operation timed out")
	case ERROR_OVERFLOW:
		return errors.New("overflow")
	case ERROR_PIPE:
		return errors.New("pipe error")
	case ERROR_INTERRUPTED:
		return errors.New("system call interrupted")
	case ERROR_NO_MEM:
		return errors.New("insufficient memory")
	case ERROR_NOT_SUPPORTED:
		return errors.New("operation not supported")
	case ERROR_INVALID_DEVICE:
		return errors.New("device is not UVC-complian")
	case ERROR_INVALID_MODE:
		return errors.New("mode not supported")
	case ERROR_CALLBACK_EXISTS:
		return errors.New("resource has a callback, can't use polling and async")
	case ERROR_OTHER:
		return errors.New("undefined error")
	default:
		return fmt.Errorf("unknown error: %d", et)
	}
}
