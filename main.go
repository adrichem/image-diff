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

const (
	statusPath     = "/status"
	diffPath       = "/"
	ignoreColorKey = "ignoreColor"
)

func init() {
	flag.StringVar(&listen, "listen", "0.0.0.0:80", "host and port to listen to")
}

func status(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func diff(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 * 1024 * 1024)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(r.MultipartForm.File) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
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

	//Decode files into Image objects
	file1, _, err := r.FormFile(files[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file1.Close()

	file2, _, err := r.FormFile(files[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file2.Close()

	img1, _, err := image.Decode(file1)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	img2, _, err := image.Decode(file2)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pxComparer := imagediff.NewYIQPixelComparer()
	//pxComparer := imagediff.NewDefaultPixelComparer()

	//If user has specified an ignore color then tell the pixelcomparer what it is
	strIgnoreColor := []byte(r.FormValue(ignoreColorKey))

	if len(strIgnoreColor) > 0 {
		ignoreColor := &color.RGBA{}
		err := json.Unmarshal(strIgnoreColor, ignoreColor)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		pxComparer.IgnoreColor = ignoreColor
	}

	//Calculate difference between the two Images
	numDiff, diffImage, err := imagediff.Diff(img1, img2, pxComparer)

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

func main() {
	handler := http.NewServeMux()
	handler.HandleFunc(statusPath, status)
	handler.HandleFunc(diffPath, diff)

	server := &http.Server{
		Addr:    listen,
		Handler: handler,
	}
	fmt.Println("Listening on " + listen)
	server.ListenAndServe()
}
