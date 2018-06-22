package imagediff

import (
	"image"
	"image/color"
)

// YIQPixelComparer uses perceptual color difference metrics to determine if 2 pixels are effectively the same.
// See http://www.progmat.uaem.mx:8080/artVol2Num2/Articulo3Vol2Num2.pdf
//
// this implementation taken from https://github.com/mapbox/pixelmatch
type YIQPixelComparer struct {
	//A color that forces both pixel to be considered equal.
	IgnoreColor *color.RGBA
	DiffColor   color.NRGBA
	Threshold   float32
}

// NewYIQPixelComparer creates a YIQPixelComparer.
func NewYIQPixelComparer() YIQPixelComparer {
	return YIQPixelComparer{
		Threshold: 0.1,
		DiffColor: color.NRGBA{R: 255, G: 0, B: 0, A: 255},
	}
}

// Compare determines whether 2 pixels are equal considering the criteria of YIQPixelComparer
func (comparer YIQPixelComparer) Compare(x int, y int, img1 image.Image, img2 image.Image) (bool, color.NRGBA) {

	// maximum acceptable square distance between two colors;
	// 35215 is the maximum possible value for the YIQ difference metric
	maxDelta := 35215 * comparer.Threshold * comparer.Threshold

	delta := colorDelta(img1, img2, x, y, false)

	// the color difference is above the threshold
	if delta > maxDelta {
		// found substantial difference not caused by anti-aliasing; draw it as the difference color
		return false, comparer.DiffColor
	}
	// pixels are similar; draw background as grayscale image blended with white
	val := blend(grayPixel(img1, x, y), 0.1)
	outputPixel := color.NRGBA{}
	outputPixel.R = uint8(val)
	outputPixel.G = uint8(val)
	outputPixel.B = uint8(val)
	outputPixel.A = 255
	return true, outputPixel

}

func color2nrgba(c color.Color) color.NRGBA {
	r, g, b, a := c.RGBA()
	return rgba2nrgba(r, g, b, a)

}

func rgba2nrgba(r, g, b, a uint32) color.NRGBA {
	return color.NRGBA{
		R: uint8(r / 0x101),
		G: uint8(g / 0x101),
		B: uint8(b / 0x101),
		A: uint8(a / 0x101),
	}
}

// calculate color difference according to the paper "Measuring perceived color difference
// using YIQ NTSC transmission color space in mobile applications" by Y. Kotsarenko and F. Ramos
func colorDelta(img1 image.Image, img2 image.Image, xPos int, yPos int, yOnly bool) float32 {

	tmp1 := img1.At(xPos, yPos)
	tmp2 := img2.At(xPos, yPos)

	Pixel1 := color2nrgba(tmp1)
	Pixel2 := color2nrgba(tmp2)

	a1 := float32(Pixel1.A) / 255
	a2 := float32(Pixel2.A) / 255

	r1 := blend(float32(Pixel1.R), a1)
	g1 := blend(float32(Pixel1.G), a1)
	b1 := blend(float32(Pixel1.B), a1)

	r2 := blend(float32(Pixel2.R), a2)
	g2 := blend(float32(Pixel2.G), a2)
	b2 := blend(float32(Pixel2.B), a2)

	y := rgb2y(r1, g1, b1) - rgb2y(r2, g2, b2)

	if yOnly {
		return y // brightness difference only
	}

	i := rgb2i(r1, g1, b1) - rgb2i(r2, g2, b2)
	q := rgb2q(r1, g1, b1) - rgb2q(r2, g2, b2)

	return 0.5053*y*y + 0.299*i*i + 0.1957*q*q
}

func rgb2y(r, g, b float32) float32 { return r*0.29889531 + g*0.58662247 + b*0.11448223 }
func rgb2i(r, g, b float32) float32 { return r*0.59597799 - g*0.27417610 - b*0.32180189 }
func rgb2q(r, g, b float32) float32 { return r*0.21147017 - g*0.52261711 + b*0.31114694 }

// blend semi-transparent color with white
func blend(c float32, a float32) float32 { return 255 + float32(c-255)*a }

func grayPixel(img image.Image, x, y int) float32 {

	tmp := img.At(x, y)
	Pixel := color2nrgba(tmp)

	a := float32(Pixel.A) / 255
	r := blend(float32(Pixel.R), a)
	g := blend(float32(Pixel.G), a)
	b := blend(float32(Pixel.B), a)
	return rgb2y(r, g, b)
}
