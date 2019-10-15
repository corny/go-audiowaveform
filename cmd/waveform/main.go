package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/corny/go-waveform"
)

var (
	blue color.Color = color.RGBA{0, 0, 255, 255}
	red  color.Color = color.RGBA{0, 0, 255, 64}
)

func main() {
	width := flag.Int("width", 500, "width in pixels")
	height := flag.Int("height", 100, "height in pixels")
	output := flag.String("output", "waveform.png", "path to output filename")
	flag.Parse()

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	wf, err := waveform.ReadWaveform(file)
	if err != nil {
		panic(err)
	}

	log.Printf("%+v", wf)

	// Draw image
	m := image.NewRGBA(image.Rect(0, 0, *width, *height))
	i := 0
	wf.EachLine(*width, func(avgMin, avgMax, peakMin, peakMax float32) {

		log.Println(peakMin, peakMax)

		draw.Draw(m, image.Rect(
			i, *height/2-int(peakMin*float32(*height)/2),
			i+1, *height/2-int(peakMax*float32(*height)/2),
		),
			&image.Uniform{red},
			image.ZP,
			draw.Src,
		)

		draw.Draw(m, image.Rect(
			i, *height/2-int(avgMin*float32(*height)/2),
			i+1, *height/2-int(avgMax*float32(*height)/2),
		),
			&image.Uniform{blue},
			image.ZP,
			draw.Src,
		)

		i++
	})

	w, err := os.Create(*output)
	if err != nil {
		panic(err)
	}
	defer w.Close()
	png.Encode(w, m)
}
