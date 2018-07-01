package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"strconv"

	"github.com/adrichem/image-diff/imagediff"
)

var listen string
var simple bool

const (
	statusPath     = "/status"
	diffSmartPath  = "/"
	diffSimplePath = "/simple"
	ignoreColorKey = "ignoreColor"
)

func init() {
	flag.StringVar(&listen, "listen", "0.0.0.0:80", "host and port to listen to")
}

func status(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func diffSmart(w http.ResponseWriter, r *http.Request) {

	img1, img2, ignoreColor, err := parseForm(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pxComparer := imagediff.NewSmartImageComparer()

	//If user has specified an ignore color then tell the pixelcomparer what it is
	if nil != ignoreColor {
		pxComparer.SetIgnoreColor(ignoreColor)
	}

	//Calculate difference between the two Images
	numDiff, diffImage, err := pxComparer.CompareImages(img1, img2)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Return the differences to the caller
	w.Header().Add("Content-Type", "image/png")
	w.Header().Add("numDifferentPixels", strconv.Itoa(numDiff))
	w.WriteHeader(http.StatusOK)
	png.Encode(w, diffImage)
	return
}

func diffSimple(w http.ResponseWriter, r *http.Request) {

	img1, img2, ignoreColor, err := parseForm(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pxComparer := imagediff.NewSimpleImageComparer()

	//If user has specified an ignore color then tell the pixelcomparer what it is
	if nil != ignoreColor {
		pxComparer.SetIgnoreColor(ignoreColor)
	}

	//Calculate difference between the two Images
	numDiff, diffImage, err := pxComparer.CompareImages(img1, img2)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	//Return the differences to the caller
	w.Header().Add("Content-Type", "image/png")
	w.Header().Add("numDifferentPixels", strconv.Itoa(numDiff))
	w.WriteHeader(http.StatusOK)
	png.Encode(w, diffImage)
	return
}

func parseForm(r *http.Request) (image.Image, image.Image, *color.NRGBA, error) {
	err := r.ParseMultipartForm(10 * 1024 * 1024)

	if err != nil {
		return nil, nil, nil, err
	}

	if len(r.MultipartForm.File) != 2 {
		return nil, nil, nil, err
	}

	var files [2]string
	var i int

	for k := range r.MultipartForm.File {
		files[i] = k
		i++
		if i >= 2 {
			break
		}
	}

	//Decode files into Image interfaces
	file1, _, err := r.FormFile(files[0])
	if err != nil {
		return nil, nil, nil, err
	}
	defer file1.Close()

	file2, _, err := r.FormFile(files[1])
	if err != nil {
		return nil, nil, nil, err
	}
	defer file2.Close()

	img1, _, err := image.Decode(file1)
	if err != nil {
		return nil, nil, nil, err
	}

	img2, _, err := image.Decode(file2)
	if err != nil {
		return nil, nil, nil, err
	}

	strIgnoreColor := []byte(r.FormValue(ignoreColorKey))
	var ignoreColor *color.NRGBA
	ignoreColor = nil
	if len(strIgnoreColor) > 0 {
		tmp := color.NRGBA{}
		err = json.Unmarshal(strIgnoreColor, &tmp)
		if err != nil {
			return nil, nil, nil, err
		}

		ignoreColor = &tmp

	}

	return img1, img2, ignoreColor, nil
}

func main() {
	handler := http.NewServeMux()
	handler.HandleFunc(statusPath, status)
	handler.HandleFunc(diffSmartPath, diffSmart)
	handler.HandleFunc(diffSimplePath, diffSimple)

	server := &http.Server{
		Addr:    listen,
		Handler: handler,
	}
	fmt.Println("Listening on " + listen)
	server.ListenAndServe()
}
