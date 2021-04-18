package uvc

// #include <libuvc/libuvc.h>
import "C"

type FrameFormat C.enum_uvc_frame_format

const (
	FRAME_FORMAT_UNKNOWN FrameFormat = C.UVC_FRAME_FORMAT_UNKNOWN
	// Any supported format
	FRAME_FORMAT_ANY          FrameFormat = C.UVC_FRAME_FORMAT_ANY
	FRAME_FORMAT_UNCOMPRESSED FrameFormat = C.UVC_FRAME_FORMAT_UNCOMPRESSED
	FRAME_FORMAT_COMPRESSED   FrameFormat = C.UVC_FRAME_FORMAT_COMPRESSED
	// YUYV/YUV2/YUV422: YUV encoding with one luminance value per pixel and
	// one UV (chrominance) pair for every two pixels.
	FRAME_FORMAT_YUYV FrameFormat = C.UVC_FRAME_FORMAT_YUYV
	FRAME_FORMAT_UYVY FrameFormat = C.UVC_FRAME_FORMAT_UYVY
	// 24-bit RGB
	FRAME_FORMAT_RGB FrameFormat = C.UVC_FRAME_FORMAT_RGB
	FRAME_FORMAT_BGR FrameFormat = C.UVC_FRAME_FORMAT_BGR
	// Motion-JPEG (or JPEG) encoded images
	FRAME_FORMAT_MJPEG FrameFormat = C.UVC_FRAME_FORMAT_MJPEG
	// Greyscale images
	FRAME_FORMAT_GRAY8  FrameFormat = C.UVC_FRAME_FORMAT_GRAY8
	FRAME_FORMAT_GRAY16 FrameFormat = C.UVC_FRAME_FORMAT_GRAY16
	// Raw colour mosaic images
	FRAME_FORMAT_BY8    FrameFormat = C.UVC_FRAME_FORMAT_BY8
	FRAME_FORMAT_BA81   FrameFormat = C.UVC_FRAME_FORMAT_BA81
	FRAME_FORMAT_SGRBG8 FrameFormat = C.UVC_FRAME_FORMAT_SGRBG8
	FRAME_FORMAT_SGBRG8 FrameFormat = C.UVC_FRAME_FORMAT_SGBRG8
	FRAME_FORMAT_SRGGB8 FrameFormat = C.UVC_FRAME_FORMAT_SRGGB8
	FRAME_FORMAT_SBGGR8 FrameFormat = C.UVC_FRAME_FORMAT_SBGGR8
	// Number of formats understood
	FRAME_FORMAT_COUNT FrameFormat = C.UVC_FRAME_FORMAT_COUNT
)

type Frame struct {
	Width  int
	Height int
}
