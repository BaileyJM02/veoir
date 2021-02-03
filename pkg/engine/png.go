package engine

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"

	"github.com/baileyjm02/veoir/pkg/queue"
	"github.com/baileyjm02/veoir/pkg/types"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"github.com/sirupsen/logrus"
)

// StartSVGQueue subscribes to the SVG event channel
func StartPNGQueue() {
	channel := make(chan queue.DataEvent)
	queue.Queues.Subscribe("engine.png", channel)
	for {
		select {

		case d := <-channel:
			go BuildPNG(d.Data.(types.Image))
		}
	}
}

// BuildPNG is a helper function to create the image that will be sent
func BuildPNG(img types.Image) {
	w, h := 512, 512
	icon, err := oksvg.ReadIconStream(strings.NewReader(img.Payload))
	if err != nil {
		logrus.Error(err)
		return
	}
	icon.SetTarget(0, 0, float64(w), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)
	out, err := os.Create(fmt.Sprintf("public/%v.png", img.Hash))
	if err != nil {
		logrus.Error(err)
		return
	}

	defer out.Close()

	err = png.Encode(out, rgba)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("SVG %v created.", img.Hash)
}
