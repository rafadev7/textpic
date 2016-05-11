package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"golang.org/x/image/math/fixed"

	"golang.org/x/image/font"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

var (
	dpi      = float64(72.0)
	fontfile = "DejaVuSansMono.ttf"
	hinting  = "full" // or "none"
	size     = float64(12)
	spacing  = float64(1)
	wonb     = false // White or Black
	f        *truetype.Font
)

func init() {

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Println("Err loading font: ", err.Error())
		return
	}
	f, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println("Err parsing font: ", err.Error())
		return
	}

}

var grayScale = "#$@B%8&WM*oahkbdpqwmZO0QLCJUYXzcvunxrjft{}[]/\\|()1?-_+~<>i!lI;:,^'''.  "

func TransformImage(url string, w, h int) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return []byte{}, errors.New("FUCK U:" + err.Error())
	}

	w2 := int(w / int(size*spacing))
	img, w2, h2 := ScaleImage(img, w2)

	text := Convert2Ascii(img, w2, h2)

	// Initialize the context.
	fg, bg := image.Black, image.White
	if wonb {
		fg, bg = image.White, image.Black
	}

	// HOW TO KNOW THE SIZE OF THE IMAGE?
	// I DUNNO...
	rgba := image.NewRGBA(image.Rect(0, 0, w2*int(size*spacing), h2*int(size*spacing*1.6))) // 640, 480
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	switch hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the text.
	pt := freetype.Pt(0, int(size)) // int(c.PointToFixed(>>6)
	var r fixed.Point26_6
	for _, s := range text {
		r, err = c.DrawString(s, pt)
		if err != nil {
			log.Println("Erro drawing:", err)
			return []byte{}, err
		}
		pt.Y += c.PointToFixed(size * spacing)
	}

	// I'm doing it because I just dunno
	// how to predict the correct image size
	cropped, err := cutter.Crop(rgba, cutter.Config{
		Width:   r.X.Round(),
		Height:  r.Y.Round(),
		Options: cutter.Copy,
	})

	// // Save that RGBA image to disk.
	// outFile, err := os.Create("out.png")
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	// defer outFile.Close()
	// bio := bufio.NewWriter(outFile)

	buf := new(bytes.Buffer)
	err = png.Encode(buf, cropped)
	if err != nil {
		log.Println(errors.New("err encoding:" + err.Error()))
		//os.Exit(1)
		return []byte{}, errors.New("Error encoding:" + err.Error())
	}
	// err = bio.Flush()
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	fmt.Println("Wrote out.png OK.")

	return buf.Bytes(), nil
}

func ScaleImage(img image.Image, w int) (image.Image, int, int) {
	sz := img.Bounds()
	h := (sz.Max.Y * w * 10) / (sz.Max.X * 16)
	img = resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	return img, w, h
}

func Convert2Ascii(img image.Image, w, h int) []string {
	text := []string{}

	for i := 0; i < h; i++ {
		line := ""
		for j := 0; j < w; j++ {
			g := color.GrayModel.Convert(img.At(j, i))
			y := reflect.ValueOf(g).FieldByName("Y").Uint()
			pos := int(y * 70 / 255)
			line = line + string(grayScale[pos])
		}
		text = append(text, line)
	}
	return text
}
