package engine

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	svg "github.com/ajstarks/svgo"
	"github.com/baileyjm02/veoir/pkg/queue"
	"github.com/baileyjm02/veoir/pkg/types"
	"github.com/sirupsen/logrus"
)

// StartSVGQueue subscribes to the SVG event channel
func StartSVGQueue() {
	channel := make(chan queue.DataEvent, 2)
	queue.Queues.Subscribe("engine.svg", channel)
	for {
		select {

		case d := <-channel:
			go BuildSVG(d.Data.(types.Image))
		}
	}
}

// BuildSVG is a helper function to create the image that will be sent
func BuildSVG(image types.Image) {
	lines := strings.Split(image.Payload, "\n")

	file, err := os.Create(fmt.Sprintf("public/%v.svg", image.Hash))
	if err != nil {
		logrus.Error(err)
		return
	}
	defer file.Close()

	if image.Theme == "light" {
		drawImage(file, lines, "#E6E6FA", "#fff", "#000")
	} else {
		drawImage(file, lines, "#a4a4a4", "#000", "#fff")
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Error(err)
		return
	}

	image.Payload = string(bytes)

	queue.Queues.Publish("engine.png", image)

	logrus.Infof("SVG %v created.", image.Hash)
}

func drawImage(w io.Writer, lines []string, bg, box, text string) {
	width := 0
	for _, line := range lines {
		length := (strings.Count(line, "\t") * 2) + len(line)
		if length > width {
			width = length
		}
	}
	width = (width * 10) + 100
	height := (len(lines) * 26) + 30
	canvas := svg.New(w)
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:"+bg)
	canvas.Roundrect(height/10, height/10, width-(2*(height/10)), height-(2*(height/10)), 12, 12, "fill:"+box)
	canvas.Circle(width-(height/10)-15, (height/10)+14, 6, "fill:tomato")
	canvas.Circle(width-(height/10)-30, (height/10)+14, 6, "fill:#FFDF00")
	canvas.Circle(width-(height/10)-45, (height/10)+14, 6, "fill:lime")
	for i, line := range lines {
		canvas.Text((height/10)+20+(strings.Count(line, "\t")*16), (height/10)+35+(i*20), line, "font-family:'JetBrains Mono';font-size:16px;fill:"+text)
	}
	canvas.End()
}
