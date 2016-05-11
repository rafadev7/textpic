package main

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"reflect"

	"github.com/nfnt/resize"
)

var grayScale = []byte("#$@B%8&WM*oahkbdpqwmZO0QLCJUYXzcvunxrjft{}[]/\\|()1?-_+~<>i!lI;:,^'''.  ")

var table = []byte(grayScale)

func TransformImage(url string) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", errors.New("FUCK U:" + err.Error())
	}

	img, w, h := ScaleImage(img, 28)

	b := Convert2Ascii(img, w, h)
	return string(b), nil
}

func ScaleImage(img image.Image, w int) (image.Image, int, int) {
	sz := img.Bounds()
	h := (sz.Max.Y * w * 10) / (sz.Max.X * 16)
	img = resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	return img, w, h
}

func Convert2Ascii(img image.Image, w, h int) []byte {
	buf := new(bytes.Buffer)

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			g := color.GrayModel.Convert(img.At(j, i))
			y := reflect.ValueOf(g).FieldByName("Y").Uint()
			pos := int(y * 70 / 255)
			_ = buf.WriteByte(table[pos])
		}
		_ = buf.WriteByte('\n')
	}
	return buf.Bytes()
}
