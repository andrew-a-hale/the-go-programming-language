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
	blackIndex = 0 // first color in palette
	greenIndex = 1 // next color in palette
	redIndex   = 2 // last color in palette
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print("invalid form")
		}
		cycles, err := strconv.ParseFloat(r.Form.Get("cycles"), 2)
		if err != nil {
			log.Print("cycle not in query or did not parse to float")
		}
		lissajous(w, cycles)
	}) // each request calls handler
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func lissajous(out io.Writer, cycles float64) {
	const (
		// cycles  = 5 // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 200   // image canvas
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 //  relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			var lineColor uint8
			if x < y {
				lineColor = greenIndex
			} else {
				lineColor = redIndex
			}
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), lineColor)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
