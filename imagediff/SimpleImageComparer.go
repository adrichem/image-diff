package imagediff

import (
	"errors"
	"image"
	"image/color"
)

// SimpleImageComparer considers pixels to be the same when their RGBA values are equal
type SimpleImageComparer struct {
	ignoreColor    *color.NRGBA
	useignoreColor bool
	DiffColor      color.NRGBA
}

//NewSimpleImageComparer creates a new SimpleImageComparer
func NewSimpleImageComparer() *SimpleImageComparer {
	return &SimpleImageComparer{
		DiffColor: color.NRGBA{R: 255, G: 0, B: 0, A: 255},
	}
}

//SetIgnoreColor tells the comparer to ignore anypixel of the specified color in either images
func (comparer *SimpleImageComparer) SetIgnoreColor(c *color.NRGBA) {
	comparer.ignoreColor = c
	comparer.useignoreColor = true
}

//CompareImages compares 2 images, returns the number of different pixels and an output image that highlights the differences
func (comparer *SimpleImageComparer) CompareImages(img1 image.Image, img2 image.Image) (int, *image.NRGBA, error) {
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
	var outputPixel color.NRGBA
	same := false
	for y := img1Bounds.Min.Y; y < img1Bounds.Max.Y; y++ {
		for x := img1Bounds.Min.X; x < img1Bounds.Max.X; x++ {

			Pixel1rgba := img1.At(x, y)
			Pixel2rgba := img2.At(x, y)
			P1NRGBA := color2nrgba(Pixel1rgba)
			P2NRGBA := color2nrgba(Pixel2rgba)

			r1 := P1NRGBA.R
			g1 := P1NRGBA.G
			b1 := P1NRGBA.B
			a1 := P1NRGBA.A

			r2 := P2NRGBA.R
			g2 := P2NRGBA.G
			b2 := P2NRGBA.B
			a2 := P2NRGBA.A

			if comparer.useignoreColor && isIgnorePixel(P1NRGBA.R, P1NRGBA.G, P1NRGBA.B, P1NRGBA.A, P2NRGBA.R, P2NRGBA.G, P2NRGBA.B, P2NRGBA.A, comparer.ignoreColor) {
				//These pixels should be ignored
				diffImage.SetNRGBA(x, y, niceOutputPixel(*comparer.ignoreColor))
				continue
			} else {
				same = r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
			}

			if !same {
				//These 2 pixels have different RGBA values
				numDifferentPixel++
				outputPixel = comparer.DiffColor
			} else {
				//These 2 pixels are exactly the same
				outputPixel = niceOutputPixel(P1NRGBA)
			}
			diffImage.SetNRGBA(x, y, outputPixel)
		}
	}

	return numDifferentPixel, diffImage, nil
}
