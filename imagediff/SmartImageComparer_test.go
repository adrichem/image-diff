package imagediff

import (
	"image"
	"image/color"
	_ "image/png"
	"strconv"
	"testing"
)

var smartComparer = NewSmartImageComparer()

func TestSmartErrorWhenImagesHaveDifferentSize(t *testing.T) {
	OnePixelBig := image.Rect(0, 0, 1, 1)
	Bigger := image.Rect(0, 0, 2, 2)

	img1 := image.NewNRGBA(OnePixelBig)
	img2 := image.NewNRGBA(Bigger)

	_, _, err := smartComparer.CompareImages(img1, img2)

	if err == nil {
		t.Fatal("Expected error when comparing different image sizes. No error was raised")
	}
}
func TestSmartNotSignificantlyDifferent(t *testing.T) {
	smartComparer.SetThreshold(0.1)
	OnePixelBig := image.Rect(0, 0, 1, 1)

	img1 := image.NewNRGBA(OnePixelBig)
	img2 := image.NewNRGBA(OnePixelBig)

	img1.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	img2.Set(0, 0, color.NRGBA{R: 255, G: 1, B: 1, A: 255})

	n, _, err := smartComparer.CompareImages(img1, img2)

	if err != nil {
		t.Fatal(err)
		return
	}

	if n != 0 {
		t.FailNow()
	}

}

func TestSmartVeryDifferent(t *testing.T) {
	smartComparer.SetThreshold(0.1)
	OnePixelBig := image.Rect(0, 0, 1, 1)

	img1 := image.NewNRGBA(OnePixelBig)
	img2 := image.NewNRGBA(OnePixelBig)

	img1.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	img2.Set(0, 0, color.NRGBA{R: 128, G: 128, B: 128, A: 255})

	n, _, err := smartComparer.CompareImages(img1, img2)

	if err != nil {
		t.Fatal(err)
		return
	}

	if n != 1 {
		t.FailNow()
	}
}
func TestSmartSame(t *testing.T) {
	smartComparer.SetThreshold(0.1)
	OnePixelBig := image.Rect(0, 0, 1, 1)

	img1 := image.NewNRGBA(OnePixelBig)
	img2 := image.NewNRGBA(OnePixelBig)

	img1.Set(0, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	img2.Set(0, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 255})

	n, _, err := smartComparer.CompareImages(img1, img2)

	if err != nil {
		t.Fatal(err)
		return
	}

	if n != 0 {
		t.FailNow()
	}
}

func TestSmartIgnoreColor(t *testing.T) {
	smartComparer.SetThreshold(0.1)
	ignoreColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	smartComparer.SetIgnoreColor(&ignoreColor)
	OnePixelBig := image.Rect(0, 0, 1, 1)

	img1 := image.NewNRGBA(OnePixelBig)
	img2 := image.NewNRGBA(OnePixelBig)

	//1 pixel is set to ignore color
	img1.Set(0, 0, ignoreColor)
	img2.Set(0, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 255})

	n, _, err := smartComparer.CompareImages(img1, img2)

	if err != nil {
		t.Fatal(err)
		return
	}

	//Other pixel is set to ignore color
	if n != 0 {
		t.Error("Expected 0 different pixels, got " + strconv.Itoa(n))
		t.FailNow()
	}

	img1.Set(0, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	img2.Set(0, 0, ignoreColor)
	n, _, err = smartComparer.CompareImages(img1, img2)

	if err != nil {
		t.Fatal(err)
		return
	}

	if n != 0 {
		t.Error("Expected 0 different pixels, got " + strconv.Itoa(n))
		t.FailNow()
	}

	//Both pixels are to to ignore color
	img1.Set(0, 0, ignoreColor)
	img2.Set(0, 0, ignoreColor)
	n, _, err = smartComparer.CompareImages(img1, img2)

	if err != nil {
		t.Fatal(err)
		return
	}

	if n != 0 {
		t.Error("Expected 0 different pixels, got " + strconv.Itoa(n))
		t.FailNow()
	}

}
