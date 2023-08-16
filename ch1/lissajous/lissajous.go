// Lissajous generates GIF animations of random Lissajous figures.
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
)

var palette = []color.Color{color.Black, color.RGBA{0, 255, 0, 255}, color.RGBA{255, 0, 0, 255}}

const (
	blackIndex uint8 = 0 // first color in palette
	greenIndex uint8 = 1 // next color in palette
	redIndex   uint8 = 2 // last color in palette
)

func main() {
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles  = 3     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 200   // image canvas
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 //  relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	lineColor := greenIndex

	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)

		if i%16 == 0 {
			lineColor = cycleColor(lineColor)
		}

		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), lineColor)
		}

		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

func cycleColor(colorIndex uint8) uint8 {
	// cycle palette avoiding black colorIndex = 0
	var newColorIndex uint8
	paletteSize := uint8(len(palette))

	newColorIndex = colorIndex + 1
	if newColorIndex % paletteSize == 0 {
		newColorIndex = 1
	}

	return newColorIndex
}