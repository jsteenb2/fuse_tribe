package tribe

import (
	"testing"
)

func TestMain(t *testing.T) {
	video := videoData{Height: 1050, Width: 1680}
	Runner(video, "test.json")
}

type videoData struct {
	Height int
	Width  int
}

func (v videoData) GetHeight() uint {
	return uint(v.Height)
}

func (v videoData) GetWidth() uint {
	return uint(v.Width)
}
