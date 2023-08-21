// Mandelbrot emits a PNG image of the Mandelbrot fractal
package main

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"math/cmplx"
	"net/http"
	"strconv"
)

const (
	xmin, ymin, xmax, ymax = -2, -2, 2, 2
	width, height          = 2048, 2048
)

func main() {
	// / endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("fractal home!"))
	})

	// /fractal endpoint
	http.HandleFunc("/fractal", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print("invalid form")
		}
		var algorithm int64
		algorithm, err := strconv.ParseInt(r.Form.Get("al"), 10, 64)
		if err != nil {
			algorithm = 1
			log.Printf("did not parse al to int, got %s, defaulted to %d", r.Form.Get("al"), algorithm)
		}
		makepng(w, r, algorithm)
	})

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func makepng(w io.Writer, r *http.Request, al int64) {
	img := image.NewRGBA(image.Rect(0, 0, width/2, height/2))
	for py := 0; py < height/2; py++ {
		for px := 0; px < width/2; px++ {
			img.Set(px, py, supersample(al, px, py))
		}
	}
	png.Encode(w, img)
}

func supersample(al int64, px, py int) color.Color {
	r1, g1, b1, a1 := compute(al, 2*px, 2*py).RGBA()
	r2, g2, b2, a2 := compute(al, 2*px, 2*py+1).RGBA()
	r3, g3, b3, a3 := compute(al, 2*px+1, 2*py).RGBA()
	r4, g4, b4, a4 := compute(al, 2*px+1, 2*py+1).RGBA()

	return color.RGBA{
		uint8((r1 + r2 + r3 + r4) / 4),
		uint8((g1 + g2 + g3 + g4) / 4),
		uint8((b1 + b2 + b3 + b4) / 4),
		uint8((a1 + a2 + a3 + a4) / 4),
	}
}

func compute(al int64, px, py int) color.Color {
	x := float64(px)/width*(xmax-xmin) + xmin
	y := float64(py)/width*(ymax-ymin) + ymin

	var color_fn func(complex128) color.Color
	switch al {
	case 1:
		color_fn = mandelbrot
	case 2:
		color_fn = julia
	default:
		log.Fatal("invalid al")
	}

	return color_fn(complex(x, y))
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast int = (1<<24 - 1) / iterations

	var v complex128
	for n := 0; n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			colorInt := contrast * n
			b := uint8(colorInt)
			g := uint8(colorInt >> 8)
			r := uint8(colorInt >> 16)
			return color.RGBA{b, g, r, 255}
		}
	}
	return color.Black
}

func julia(z complex128) color.Color {
	const iterations = 200
	roots := [4]complex128{complex(1, 0), complex(-1, 0), complex(0, 1), complex(0, -1)}

	var v complex128 = z
	var c color.Color = color.Black
	for n := 0; n < iterations; n++ {
		v = v - (v*v*v*v-1)/(4*v*v*v)
		switch {
		case math.IsNaN(real(v)) || math.IsNaN(imag(v)):
			return c
		case near(v, roots[0]):
			return color.RGBA{228, 26, 28, uint8(n * 100)}
		case near(v, roots[1]):
			return color.RGBA{55, 126, 184, uint8(n * 100)}
		case near(v, roots[2]):
			return color.RGBA{77, 175, 74, uint8(n * 100)}
		case near(v, roots[3]):
			return color.RGBA{152, 78, 163, uint8(n * 100)}
		}
	}
	return c
}

func near(v1, v2 complex128) bool {
	const tol float64 = 0.0001
	if cmplx.Abs(v1-v2) < tol {
		return true
	}
	return false
}
