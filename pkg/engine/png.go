package engine

import (
	"fmt"

	"github.com/baileyjm02/veoir/pkg/queue"
	"github.com/baileyjm02/veoir/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/tdewolff/canvas/rasterizer"
)

// StartPNGQueue subscribes to the PNG event channel
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
func BuildPNG(image types.Image) {
	// Output the rasterized PNG to the public/<hash>.png file
	err := image.Canvas.WriteFile(fmt.Sprintf("public/%v.png", image.Hash), rasterizer.PNGWriter(2.5))
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("PNG %v created.", image.Hash)
}
