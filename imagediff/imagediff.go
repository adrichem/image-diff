package imagediff

import (
	"errors"
	"image"
)

// Diff determines the differences between 2 images of the same size.
// Every pixel is compared using the supplied PixelComparer and the number of different pixels and an output image
// is returned. The colors of different and same pixels are decided by the PixelComparer
func Diff(img1 image.Image, img2 image.Image, comparer PixelComparer) (int, *image.NRGBA, error) {

	//Size must be same
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

	for y := img1Bounds.Min.Y; y < img1Bounds.Max.Y; y++ {
		for x := img1Bounds.Min.X; x < img1Bounds.Max.X; x++ {
			same, outputPixel := comparer.Compare(x, y, img1, img2)
			if !same {
				numDifferentPixel++
			}
			diffImage.SetNRGBA(x, y, outputPixel)
		}
	}

	return numDifferentPixel, diffImage, nil
}
