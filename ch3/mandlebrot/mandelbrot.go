// Mandelbrot emits a PNG image of the Mandelbrot fractal
package main

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math/cmplx"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		makepng(w, r)
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func makepng(w io.Writer, r *http.Request) {
	const (
		xmin, ymin, xmax, ymax = -2, -2, 2, 2
		width, height          = 4096, 4096
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			// Image point (px, py) represents complex value z.
			pz := mandelbrot(z)
			img.Set(px, py, pz)
			img.Set(px, py+1, pz)
			img.Set(px, py+2, pz)
			img.Set(px, py+3, pz)
			img.Set(px+1, py+1, pz)
			img.Set(px+1, py+2, pz)
			img.Set(px+1, py+3, pz)
			img.Set(px+2, py+1, pz)
			img.Set(px+2, py+2, pz)
			img.Set(px+2, py+3, pz)
			img.Set(px+3, py+1, pz)
			img.Set(px+3, py+2, pz)
			img.Set(px+3, py+3, pz)
			px += 4
		}
		py += 4
	}
	png.Encode(w, img) // NOTE: ignoring errors
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const bitlength = 1<<8 - 1
	var contrast = (1<<24 - 1) / iterations

	var v complex128
	for n := 0; n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			colorInt := contrast * n
			b := uint8(colorInt)
			g := uint8(colorInt >> 8)
			r := uint8(colorInt >> 16)
			return color.RGBA{r, g, b, 255}
		}
	}
	return color.Black
}
