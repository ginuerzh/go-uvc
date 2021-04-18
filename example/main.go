package main

import (
	"log"

	gouvc "github.com/ginuerzh/go-uvc"
)

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	uvc := &gouvc.UVC{}
	if err := uvc.Init(); err != nil {
		log.Fatal("init: ", err)
	}
	defer uvc.Exit()
	log.Println("UVC initialized")

	dev, err := uvc.FindDevice(0, 0, "")
	if err != nil {
		log.Fatal("find device: ", err)
	}
	defer dev.Unref()
	log.Println("device found")

	if err := dev.Open(); err != nil {
		log.Fatal(err)
	}
	defer dev.Close()
	log.Println("device opened")

	dev.Diag()

	info, err := dev.Info()
	log.Printf("find device: %04x:%04x (%s, %s)",
		info.VendorID, info.ProductID, info.Manufacturer, info.Product)

	if err := dev.SetAEMode(gouvc.AEModeManual); err != nil {
		log.Fatal("set ae mode:", err)
	}

	var frameFormat gouvc.FrameFormat
	formatDesc := dev.GetFormatDesc()
	switch formatDesc.Subtype {
	case gouvc.VS_FORMAT_MJPEG:
		frameFormat = gouvc.FRAME_FORMAT_MJPEG
	case gouvc.VS_FORMAT_FRAME_BASED:
		// frameFormat = gouvc.FRAME_FORMAT_H264
	default:
		frameFormat = gouvc.FRAME_FORMAT_YUYV
	}

	width := 640
	height := 480
	fps := 30

	if frameDesc := formatDesc.FrameDesc(); frameDesc != nil {
		width = int(frameDesc.Width)
		height = int(frameDesc.Height)
		fps = int(10000000 / frameDesc.DefaultFrameInterval)
	}

	log.Printf("First format: %dx%d %dfps", width, height, fps)

	stream, err := dev.GetStream(frameFormat, width, height, fps)
	if err != nil {
		log.Fatal("get stream:", err)
	}

	ctrl := stream.Ctrl()
	log.Println(ctrl.String())

	if err := stream.Open(); err != nil {
		log.Fatal("open stream:", err)
	}
	defer stream.Close()

	cf, err := stream.Start()
	if err != nil {
		log.Fatal("start streaming:", err)
	}
	defer stream.Stop()

	for frame := range cf {
		log.Println("got image:", frame.Width, frame.Height)
	}
}
