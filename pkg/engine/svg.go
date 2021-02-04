package engine

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	"github.com/baileyjm02/veoir/pkg/queue"
	"github.com/baileyjm02/veoir/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/svg"
)

var (
	builtCanvas *canvas.Canvas
	fontFamily  *canvas.FontFamily
	nicePurple  = color.RGBA{0xE6, 0xE6, 0xFA, 0xff}
	niceBlue    = color.RGBA{0x49, 0x46, 0x80, 0xff}
	tomato      = color.RGBA{0xFF, 0x63, 0x47, 0xff}
	yellow      = color.RGBA{0xE6, 0xDF, 0x00, 0xff}
	green       = color.RGBA{0x94, 0xFF, 0x00, 0xff}
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

	fontFamily = canvas.NewFontFamily("times")
	fontFamily.Use(canvas.CommonLigatures)
	if err := fontFamily.LoadFontFile("JetBrainsMono.ttf", canvas.FontRegular); err != nil {
		panic(err)
	}

	if image.Theme == "light" {
		builtCanvas = drawImage(lines, nicePurple, color.White, color.Black)
	} else {
		builtCanvas = drawImage(lines, niceBlue, color.Black, color.White)
	}

	image.Canvas = *builtCanvas
	go queue.Queues.Publish("engine.png", image)

	err = svg.Writer(file, builtCanvas)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("SVG %v created.", image.Hash)
}

func drawImage(lines []string, bg, box, text color.Color) *canvas.Canvas {
	intWidth := 0
	for _, line := range lines {
		length := (strings.Count(line, "\t") * 6) + len(line)
		if length > intWidth {
			intWidth = length
		}
	}
	width := float64((intWidth * 6) + 100)
	height := float64((len(lines) * 17) + 45)

	ctx := canvas.New(float64(width), float64(height))
	c := canvas.NewContext(ctx)

	face := fontFamily.Face(30.0, text, canvas.FontRegular, canvas.FontNormal)
	rect := canvas.Rectangle(width+10, height+10)
	inner := canvas.RoundedRectangle(width-(2*(height/10)), height-(2*(height/10)), 12.0)
	circle := canvas.Circle(6.0)
	c.SetFillColor(bg)
	c.DrawPath(0, 0, rect)
	c.SetFillColor(box)
	c.DrawPath(height/10, height/10, inner)
	c.SetFillColor(tomato)
	c.DrawPath(width-(height/10)-12, height-(height/10)-12, circle)
	c.SetFillColor(yellow)
	c.DrawPath(width-(height/10)-28, height-(height/10)-12, circle)
	c.SetFillColor(green)
	c.DrawPath(width-(height/10)-44, height-(height/10)-12, circle)
	for i, line := range lines {
		_ = i
		x := (height / 10.0) + 10.0 + (float64(strings.Count(line, "\t")) * 8.0)
		y := height - (height / 10) - 35.0 - (float64(i) * 13.7)
		text := canvas.NewTextLine(face, strings.ReplaceAll(line, "\t", ""), canvas.Left)
		c.DrawText(float64(x), float64(y), text)
	}

	return ctx
}
