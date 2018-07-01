package imagediff

import (
	"errors"
	"image"
	"image/color"
)

// SmartImageComparer uses perceptual color difference metrics to determine if pixels look the same.
// See http://www.progmat.uaem.mx:8080/artVol2Num2/Articulo3Vol2Num2.pdf
type SmartImageComparer struct {
	ignoreColor    *color.NRGBA
	useignoreColor bool
	DiffColor      color.NRGBA
	threshold      float32
}

// NewSmartImageComparer constructs a SmartImageComparer.
func NewSmartImageComparer() *SmartImageComparer {
	return &SmartImageComparer{
		threshold: 0.1,
		DiffColor: color.NRGBA{R: 255, G: 0, B: 0, A: 255},
	}
}

//SetIgnoreColor tells the comparer to ignore any pixel of the specified color in either images.
func (comparer *SmartImageComparer) SetIgnoreColor(c *color.NRGBA) {
	comparer.ignoreColor = c
	comparer.useignoreColor = true
}

//SetThreshold tells the comparer how different a pixel must be before its significant.
func (comparer *SmartImageComparer) SetThreshold(t float32) {
	comparer.threshold = t
}

//CompareImages compares 2 images, returns the number of different pixels and an output image that highlights the differences.
func (comparer SmartImageComparer) CompareImages(img1 image.Image, img2 image.Image) (int, *image.NRGBA, error) {
	img1Bounds := img1.Bounds()
	img2Bounds := img2.Bounds()
	XBoundsDifferent := img1Bounds.Max.X-img1Bounds.Min.X != img2Bounds.Max.X-img2Bounds.Min.X
	YBoundsDifferent := img1Bounds.Max.Y-img1Bounds.Min.Y != img2Bounds.Max.Y-img2Bounds.Min.Y
	if XBoundsDifferent || YBoundsDifferent {
		return 0, nil, errors.New("Images not same size")
	}

	//Determine if each pixel is same or different
	numDifferentPixel := 0
	diffImage := image.NewNRGBA(img1Bounds)

	// maximum acceptable square distance between two colors;
	// 35215 is the maximum possible value for the YIQ difference metric
	maxDifference := 35215 * comparer.threshold * comparer.threshold

	for y := img1Bounds.Min.Y; y < img1Bounds.Max.Y; y++ {
		for x := img1Bounds.Min.X; x < img1Bounds.Max.X; x++ {

			Pixel1 := img1.At(x, y)
			Pixel2 := img2.At(x, y)

			P1NRGBA := color2nrgba(Pixel1)
			P2NRGBA := color2nrgba(Pixel2)

			val := blend(grayPixel(&P1NRGBA), 0.1)
			outputPixelWhenTheresNoDifference := color.NRGBA{}
			outputPixelWhenTheresNoDifference.R = uint8(val)
			outputPixelWhenTheresNoDifference.G = uint8(val)
			outputPixelWhenTheresNoDifference.B = uint8(val)
			outputPixelWhenTheresNoDifference.A = 255

			if comparer.useignoreColor && isIgnorePixel(P1NRGBA.R, P1NRGBA.G, P1NRGBA.B, P1NRGBA.A, P2NRGBA.R, P2NRGBA.G, P2NRGBA.B, P2NRGBA.A, comparer.ignoreColor) {
				diffImage.SetNRGBA(x, y, outputPixelWhenTheresNoDifference)
				continue
			}

			difference := perceptualColorDifference(&P1NRGBA, &P2NRGBA)

			if difference > maxDifference {
				numDifferentPixel++
				diffImage.SetNRGBA(x, y, comparer.DiffColor)
			} else {
				diffImage.SetNRGBA(x, y, outputPixelWhenTheresNoDifference)
			}
		}
	}

	return numDifferentPixel, diffImage, nil
}

// calculate color difference according to the paper "Measuring perceived color difference
// using YIQ NTSC transmission color space in mobile applications" by Y. Kotsarenko and F. Ramos
func perceptualColorDifference(Pixel1 *color.NRGBA, Pixel2 *color.NRGBA) float32 {

	a1 := float32(Pixel1.A) / 255
	a2 := float32(Pixel2.A) / 255

	r1 := blend(float32(Pixel1.R), a1)
	g1 := blend(float32(Pixel1.G), a1)
	b1 := blend(float32(Pixel1.B), a1)

	r2 := blend(float32(Pixel2.R), a2)
	g2 := blend(float32(Pixel2.G), a2)
	b2 := blend(float32(Pixel2.B), a2)

	y := rgb2y(r1, g1, b1) - rgb2y(r2, g2, b2)
	i := rgb2i(r1, g1, b1) - rgb2i(r2, g2, b2)
	q := rgb2q(r1, g1, b1) - rgb2q(r2, g2, b2)

	return 0.5053*y*y + 0.299*i*i + 0.1957*q*q
}

func rgb2y(r, g, b float32) float32 { return r*0.29889531 + g*0.58662247 + b*0.11448223 }
func rgb2i(r, g, b float32) float32 { return r*0.59597799 - g*0.27417610 - b*0.32180189 }
func rgb2q(r, g, b float32) float32 { return r*0.21147017 - g*0.52261711 + b*0.31114694 }

// blend semi-transparent color with white
func blend(c float32, a float32) float32 { return 255 + float32(c-255)*a }

func grayPixel(Pixel *color.NRGBA) float32 {
	a := float32(Pixel.A) / 255
	r := blend(float32(Pixel.R), a)
	g := blend(float32(Pixel.G), a)
	b := blend(float32(Pixel.B), a)
	return rgb2y(r, g, b)
}
