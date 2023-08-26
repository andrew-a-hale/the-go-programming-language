// Mandelbrot emits a PNG image of the Mandelbrot fractal
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math/big"
	"net/http"
)

const (
	xmin, ymin, xmax, ymax = -2, -2, 2, 2
	width, height          = 1024, 1024
	iterations             = 20
	contrast               = (1<<24 - 1) / iterations
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		makepng(w, r)
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func makepng(w io.Writer, r *http.Request) {
	two := big.NewRat(2, 1)
	four := big.NewRat(4, 1) // squared distance
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	var zx, zy big.Rat
	for py := 0; py < height; py++ {
		zy.SetFloat64(float64(py)/width*(ymax-ymin) + ymin)
		for px := 0; px < width; px++ {
			zx.SetFloat64(float64(px)/width*(xmax-xmin) + xmin)
			var vx, vy big.Rat
			var fill color.Color = color.Black
			for n := 0; n < iterations; n++ {
				var v1, v2, v3, v4, v5 big.Rat
				// do mandelbrot via z^2+c = x^2-y^2+2xyi+c => x*x-y*y+real(c), 2xy+imag(c)
				// real term
				v1.Mul(&vx, &vx)
				v2.Mul(&vy, &vy)
				v3.Sub(&v1, &v2)

				// imag term
				v4.Mul(&vx, &vy)
				v5.Mul(two, &v4)
				
				// update terms
				vx.Add(&v3, &zx)
				vy.Add(&v5, &zy)

				// compute distance
				v1.Mul(&vx, &vx)
				v2.Mul(&vy, &vy)
				v3.Add(&v1, &v2)

				// check is distance > 2 via pythagoras (x*x + y*y > 4)
				if v3.Cmp(four) > 0 {
					fmt.Printf("(x, y): (%d, %d) -- iter: %d\n", px, py, n)
					colorInt := contrast * n
					b := uint8(colorInt)
					g := uint8(colorInt >> 8)
					r := uint8(colorInt >> 16)
					fill = color.RGBA{b, g, r, 255}
					break
				}
			}
			img.Set(px, py, fill)
		}
	}
	png.Encode(w, img)
}