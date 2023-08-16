// Lissajous generates GIF animations of random Lissajous figures.
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

var palette = []color.Color{color.Black, color.RGBA{0, 255, 0, 255}, color.RGBA{255, 0, 0, 255}}

const (
	blackIndex uint8 = 0 // first color in palette
	greenIndex uint8 = 1 // next color in palette
	redIndex   uint8 = 2 // last color in palette
)

func main() {
	// / endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print("invalid form")
		}
		var resp []byte = []byte("a good server")
		w.Write(resp)
	})

	// /gif endpoint
	http.HandleFunc("/gif", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print("invalid form")
		}

		cycles, err := strconv.ParseFloat(r.Form.Get("cycles"), 64)
		if err != nil {
			log.Fatal("cycle not in query or did not parse to float")
		}

		lissajous(w, cycles)
	})

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func lissajous(out io.Writer, cycles float64) {
	const (
		res     = 0.001 // angular resolution
		size    = 200   // image canvas
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	_cycles := cycles            // number of complete x oscillator revolutions
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

		for t := 0.0; t < _cycles*2*math.Pi; t += res {
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

	paletteSize := uint8(len(palette) - 1)
	if colorIndex >= paletteSize {
		newColorIndex = 1
	} else {
		newColorIndex = colorIndex + 1
	}

	return newColorIndex
}
