package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/adrichem/image-diff/imagediff"
)

var (
	printUsage      bool
	algorithm       string
	listen          string
	img1            string
	img2            string
	output          string
	flagIgnoreColor bool
	ignoreR         string
	ignoreG         string
	ignoreB         string
	ignoreA         string
	comparer        imagediff.ImageComparer
)

const (
	statusPath     = "/status"
	diffPath       = "/"
	ignoreColorKey = "ignoreColor"
	smart          = "smart"
	simple         = "simple"
)

func init() {
	ignoreMessage := "Ignore differences when one of the pixels matches the specified RGBA values. If either R,G or B are not supplied, then differences will not be ignored"
	flag.BoolVar(&printUsage, "help", false, "print usage")
	flag.BoolVar(&printUsage, "h", false, "print usage")
	flag.StringVar(&listen, "listen", "0.0.0.0:80", "host and port to listen to")
	flag.StringVar(&algorithm, "algorithm", "smart", "Which algorithm to use. [smart | simple]")
	flag.StringVar(&img1, "img1", "", "Image to compare")
	flag.StringVar(&img2, "img2", "", "Other image to compare with")
	flag.StringVar(&output, "output", "", "Where to store result image")
	flag.StringVar(&ignoreR, "ignoreR", "", ignoreMessage)
	flag.StringVar(&ignoreG, "ignoreG", "", ignoreMessage)
	flag.StringVar(&ignoreB, "ignoreB", "", ignoreMessage)
	flag.StringVar(&ignoreA, "ignoreA", "255", ignoreMessage)
	flag.Parse()
}

func status(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func diffHandler(w http.ResponseWriter, r *http.Request) {

	img1, img2, ignoreColor, err := parseForm(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//If user has specified an ignore color then tell the comparer what it is
	if nil != ignoreColor {
		comparer.SetIgnoreColor(ignoreColor)
	}

	//Calculate difference between the two Images
	numDiff, diffImage, err := comparer.CompareImages(img1, img2)
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
		return nil, nil, nil, errors.New("expected 2 form files")
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

	if printUsage {
		flag.PrintDefaults()
		return
	}

	switch algorithm {
	case smart:
		comparer = imagediff.NewSmartImageComparer()
		break
	case simple:
		comparer = imagediff.NewSimpleImageComparer()
		break
	default:
		log.Fatal("Unknown algorithm: " + algorithm)
		break
	}

	runAsService := len(img1) == 0 || len(img2) == 0 || len(output) == 0
	if runAsService {
		handler := http.NewServeMux()
		handler.HandleFunc(statusPath, status)
		handler.HandleFunc(diffPath, diffHandler)

		server := &http.Server{
			Addr:    listen,
			Handler: handler,
		}
		fmt.Println("Listening on " + listen)
		log.Fatal(server.ListenAndServe())
		return
	}

	//run as stand alone executable
	img1AbsPath, _ := filepath.Abs(img1)
	file1, err := os.Open(img1AbsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file1.Close()

	img2AbsPath, _ := filepath.Abs(img2)
	file2, err := os.Open(img2AbsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	if len(ignoreR) > 0 && len(ignoreG) > 0 && len(ignoreB) > 0 {
		ignoreColor := color.NRGBA{}
		var tmpR, tmpG, tmpB, tmpA uint64

		tmpR, err := strconv.ParseUint(ignoreR, 10, 8)
		if err != nil {
			log.Fatal(err)
		}
		tmpG, err = strconv.ParseUint(ignoreG, 10, 8)
		if err != nil {
			log.Fatal(err)
		}
		tmpB, err = strconv.ParseUint(ignoreB, 10, 8)
		if err != nil {
			log.Fatal(err)
		}
		tmpA, err = strconv.ParseUint(ignoreA, 10, 8)
		if err != nil {
			log.Fatal(err)
		}

		ignoreColor.R = uint8(tmpR)
		ignoreColor.G = uint8(tmpG)
		ignoreColor.B = uint8(tmpB)
		ignoreColor.A = uint8(tmpA)

		comparer.SetIgnoreColor(&ignoreColor)
	}

	decodedImg1, _, err := image.Decode(file1)
	if err != nil {
		log.Fatal(err)
	}

	decodedImg2, _, err := image.Decode(file2)
	if err != nil {
		log.Fatal(err)
	}

	n, diffImage, err := comparer.CompareImages(decodedImg1, decodedImg2)

	if err != nil {
		log.Fatal(err)
	}

	fileOutput, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	defer fileOutput.Close()

	fmt.Println(img1 + "," + img2 + "," + strconv.Itoa(n) + "," + output)

	err = png.Encode(fileOutput, diffImage)
	if err != nil {
		log.Fatal(err)
	}

}
