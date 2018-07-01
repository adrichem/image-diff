package imagediff

import (
	"image"
	"image/color"
	_ "image/png"
	"strconv"
	"testing"
)

var pxComparer = NewSimpleImageComparer()

func TestSimpleErrorWhenImagesHaveDifferentSize(t *testing.T) {
	OnePixelBig := image.Rect(0, 0, 1, 1)
	Bigger := image.Rect(0, 0, 2, 2)

	img1 := image.NewNRGBA(OnePixelBig)
	img2 := image.NewNRGBA(Bigger)
	_, _, err := pxComparer.CompareImages(img1, img2)
	if err == nil {
		t.Fatal("Expected error when comparing different image sizes. No error was raised")
	}
}
func TestSimpleSame(t *testing.T) {
	OnePixelBig := image.Rect(0, 0, 1, 1)

	img1 := image.NewNRGBA(OnePixelBig)
	img2 := image.NewNRGBA(OnePixelBig)
	img1.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	img2.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})

	n, _, err := pxComparer.CompareImages(img1, img2)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatal("Expected 0, got " + strconv.Itoa(n))
	}
}

func TestSimpleDifferent(t *testing.T) {

	OnePixelBig := image.Rect(0, 0, 1, 1)

	img1 := image.NewNRGBA(OnePixelBig)
	img2 := image.NewNRGBA(OnePixelBig)
	img1.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	img2.Set(0, 0, color.NRGBA{R: 244, G: 1, B: 0, A: 255})

	n, _, err := pxComparer.CompareImages(img1, img2)
	if err != nil {
		t.Fatal(err)
	}

	if n != 1 {
		t.Fatal("Expected 1, got " + strconv.Itoa(n))
	}
}

func TestSimpleIgnoreColor(t *testing.T) {

	ignoreColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	pxComparer.SetIgnoreColor(&ignoreColor)
	OnePixelBig := image.Rect(0, 0, 1, 1)

	img1 := image.NewNRGBA(OnePixelBig)
	img2 := image.NewNRGBA(OnePixelBig)

	//1 pixel is set to ignore color
	img1.Set(0, 0, ignoreColor)
	img2.Set(0, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 255})

	n, _, err := pxComparer.CompareImages(img1, img2)

	if err != nil {
		t.Fatal(err)
		return
	}

	//Other pixel is set to ignore color
	if n != 0 {
		t.Fatal("Expected 0 different pixels, got " + strconv.Itoa(n))
	}

	img1.Set(0, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	img2.Set(0, 0, ignoreColor)
	n, _, err = pxComparer.CompareImages(img1, img2)

	if err != nil {
		t.Fatal(err)
	}

	if n != 0 {
		t.Fatal("Expected 0 different pixels, got " + strconv.Itoa(n))
	}

	//Both pixels are to to ignore color
	img1.Set(0, 0, ignoreColor)
	img2.Set(0, 0, ignoreColor)
	n, _, err = pxComparer.CompareImages(img1, img2)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatal("Expected 0 different pixels, got " + strconv.Itoa(n))
	}

}
