package imagediff

import (
	"image"
	"image/color"
)

// PixelComparer determines if 2 pixels are the same and returns an output color to represent result of comparing the 2 pixels
type PixelComparer interface {
	Compare(x int, y int, img1 image.Image, img2 image.Image) (bool, color.NRGBA)
}
