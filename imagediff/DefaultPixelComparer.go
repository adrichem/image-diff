package imagediff

import (
	"image"
	"image/color"
)

// DefaultPixelComparer considers 2 pixels to be equal when their RGBA values are equal or one of the pixels is equal
// to the optional ignore color
type DefaultPixelComparer struct {
	//An color that forces both pixel to be considered equal.
	IgnoreColor *color.RGBA
	DiffColor   color.NRGBA
}

//NewDefaultPixelComparer creates a new DefaultPixelComparer
func NewDefaultPixelComparer() DefaultPixelComparer {
	return DefaultPixelComparer{
		DiffColor: color.NRGBA{R: 255, G: 0, B: 0, A: 255},
	}
}

// Compare determines whether 2 pixels are equal considering the criteria of DefaultPixelComparer
func (c DefaultPixelComparer) Compare(x int, y int, img1 image.Image, img2 image.Image) (bool, color.NRGBA) {
	Pixel1 := img1.At(x, y)
	Pixel2 := img2.At(x, y)

	r1, g1, b1, a1 := Pixel1.RGBA()
	r2, g2, b2, a2 := Pixel2.RGBA()

	greyScalePixel := rgba2nrgba(color.GrayModel.Convert(Pixel1).RGBA())

	if nil != c.IgnoreColor {
		rignore, gignore, bignore, aignore := c.IgnoreColor.RGBA()
		if r1 == rignore && g1 == gignore && b1 == bignore && a1 == aignore {
			return true, greyScalePixel
		}

		if r2 == rignore && g2 == gignore && b2 == bignore && a2 == aignore {
			return true, greyScalePixel
		}
	}

	same := r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2

	if same {
		return true, greyScalePixel
	}
	return false, c.DiffColor
}
