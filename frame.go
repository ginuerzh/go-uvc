package uvc

// #include <libuvc-binding.h>
import "C"
import (
	"bytes"
	"io"
	"log"
	"time"
	"unsafe"

	"github.com/mattn/go-pointer"
)

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

// Frame is an image frame received from the UVC device.
// It implements io.Reader.
type Frame struct {
	// Width of image in pixels
	Width int
	// Height of image in pixels
	Height int
	// Pixel data format
	FrameFormat FrameFormat
	// Number of bytes per horizontal line (undefined for compressed format)
	Step int
	// Frame number (may skip, but is strictly monotonically increasing)
	Sequence uint32
	// Estimate of system time when the device started capturing the image
	CaptureTime time.Time
	// Is the data buffer owned by the library?
	// If true, the data buffer can be arbitrarily reallocated by frame conversion functions.
	// If false, the data buffer will not be reallocated or freed by the library.
	LibraryOwned bool
	// Image data for this frame
	data io.Reader
	//  Metadata for this frame if available
	Metadata []byte

	frame *C.struct_uvc_frame
}

func (fr *Frame) Read(b []byte) (int, error) {
	return fr.data.Read(b)
}

//export go_frame_cb
func go_frame_cb(frame *C.struct_uvc_frame, p unsafe.Pointer) {
	fc := pointer.Restore(p).(chan *Frame)

	switch FrameFormat(frame.frame_format) {
	case FRAME_FORMAT_YUYV:
		bgr := C.uvc_allocate_frame(C.size_t(frame.width * frame.height * 3))
		if bgr == nil {
			log.Println("unable to allocate bgr frame")
			return
		}
		defer C.uvc_free_frame(bgr)

		r := C.uvc_any2bgr(frame, bgr)
		if err := newError(ErrorType(r)); err != nil {
			log.Println(err)
			return
		}
		frame = bgr
	}

	fr := &Frame{
		Width:        int(frame.width),
		Height:       int(frame.height),
		FrameFormat:  FrameFormat(frame.frame_format),
		Step:         int(frame.step),
		Sequence:     uint32(frame.sequence),
		LibraryOwned: frame.library_owns_data > 0,
		CaptureTime:  time.Unix(int64(frame.capture_time.tv_sec), int64(frame.capture_time.tv_usec)*1000),
		data:         bytes.NewReader(C.GoBytes(unsafe.Pointer(frame.data), C.int(frame.data_bytes))),
		// Metadata:     C.GoBytes(unsafe.Pointer(frame.metadata), C.int(frame.metadata_bytes)),
		frame: frame,
	}

	select {
	case fc <- fr:
	default:
		log.Printf("frame %d dropped", frame.sequence)
	}
}
