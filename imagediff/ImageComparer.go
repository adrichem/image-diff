package imagediff

import (
	"image"
	"image/color"
)

//ImageComparer determines if images are the same. It returns the number of
//different pixels and an image highlighting the differences
type ImageComparer interface {
	CompareImages(img1 image.Image, img2 image.Image) (int, *image.NRGBA, error)
	SetIgnoreColor(IgnoreColor *color.NRGBA)
}

func isIgnorePixel(r1, g1, b1, a1, r2, g2, b2, a2 uint8, ignoreColor *color.NRGBA) bool {
	rignore := ignoreColor.R
	gignore := ignoreColor.G
	bignore := ignoreColor.B
	aignore := ignoreColor.A
	return r1 == rignore && g1 == gignore && b1 == bignore && a1 == aignore || r2 == rignore && g2 == gignore && b2 == bignore && a2 == aignore
}

func color2nrgba(c color.Color) color.NRGBA {
	r, g, b, a := c.RGBA()
	return rgba2nrgba(r, g, b, a)
}

// rgba2nrgba converts alpha-premultiplied red, green, blue and alpha values to non multplied values
func rgba2nrgba(r, g, b, a uint32) color.NRGBA {
	return color.NRGBA{
		R: uint8(r / 0x101),
		G: uint8(g / 0x101),
		B: uint8(b / 0x101),
		A: uint8(a / 0x101),
	}
}

func niceOutputPixel(p color.NRGBA) color.NRGBA {
	val := blend(grayPixel(&p), 0.1)
	output := color.NRGBA{}
	output.R = uint8(val)
	output.G = uint8(val)
	output.B = uint8(val)
	output.A = 255

	return output

}
