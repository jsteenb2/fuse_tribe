package heater

import (
	"strconv"
	"sync"
	"tribe/parser"

	"gopkg.in/gographics/imagick.v2/imagick"
)

// CreateHeatMaps creates a series of png images with coordinates from viewers gaze, fixation, and the the Video details
func CreateHeatMaps(trackerData []parser.View, video VideoInterface) {
	imagick.Initialize()
	defer imagick.Terminate()

	var wg sync.WaitGroup
	concurrency := 15
	sem := make(chan bool, concurrency)

	for i, tracker := range trackerData {
		wg.Add(1)
		sem <- true
		go func(i int64, track parser.View, video VideoInterface) {
			defer func() {
				<-sem
				wg.Done()
			}()
			CreateHeatMap(strconv.FormatInt(i, 10), track, video)
		}(int64(i), tracker, video)
	}

	wg.Wait()
}

// CreateHeatMap generates a single png image that maps the viewers gaze to the screen
func CreateHeatMap(filename string, coords ViewerInterface, video VideoInterface) {
	dw := buildDrawing(coords)
	mw := buildImage(video)
	mw.DrawImage(dw)
	mw.WriteImage("/Users/jonathansteenbergen/go/src/tribe/heater/imgs/" + filename + ".png")
}

func buildDrawing(coords ViewerInterface) *imagick.DrawingWand {
	radius := 35.0
	dw := imagick.NewDrawingWand()
	dw.PushDrawingWand()
	dw.SetStrokeColor(setColor(coords))
	dw.SetStrokeWidth(1)
	dw.SetStrokeAntialias(true)
	dw.SetStrokeLineCap(imagick.LINE_CAP_ROUND)
	dw.SetStrokeLineJoin(imagick.LINE_JOIN_ROUND)
	dw.SetFillColor(setColor(coords))
	dw.RoundRectangle(coords.GetXCoord(), coords.GetYCoord(), 2*radius+coords.GetXCoord(), 2*radius+coords.GetYCoord(), 2*radius, 2*radius)

	return dw
}

func buildImage(video VideoInterface) *imagick.MagickWand {
	mw := imagick.NewMagickWand()
	cw := imagick.NewPixelWand()
	cw.SetColor("transparent")
	mw.NewImage(video.GetWidth(), video.GetHeight(), cw)
	return mw
}

func setColor(coords ViewerInterface) *imagick.PixelWand {
	cw := imagick.NewPixelWand()
	if coords.IsFixed() {
		cw.SetColor("rgb(255,0,0)")
	} else {
		cw.SetColor("rgb(147, 21, 136)")
	}
	cw.SetAlpha(0.5)
	return cw
}

// ViewerInterface is struct with info about an image file
type ViewerInterface interface {
	GetXCoord() float64
	GetYCoord() float64
	IsFixed() bool
}

// VideoInterface is the interface that maps to an input video file
type VideoInterface interface {
	GetHeight() uint
	GetWidth() uint
}
