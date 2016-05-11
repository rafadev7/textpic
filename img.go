package main

import (
	"bufio"
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
	"os"
	"reflect"

	"golang.org/x/image/font"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
)

var (
	dpi      = float64(72.0)
	fontfile = "DejaVuSansMono.ttf"
	hinting  = "none" // or "full"
	size     = float64(13)
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

	text := Convert2Ascii(img, w, h)

	// Initialize the context.
	fg, bg := image.Black, image.White
	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	if wonb {
		fg, bg = image.White, image.Black
		ruler = color.RGBA{0x22, 0x22, 0x22, 0xff}
	}
	rgba := image.NewRGBA(image.Rect(0, 0, 640, 480))
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

	// Draw the guidelines.
	for i := 0; i < 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, 10, ruler)
	}

	// Draw the text.
	pt := freetype.Pt(10, 10+int(c.PointToFixed(size)>>6))
	for _, s := range text {
		_, err = c.DrawString(s, pt)
		if err != nil {
			log.Println("Erro drawing:", err)
			return "", nil
		}
		pt.Y += c.PointToFixed(size * spacing)
	}

	// Save that RGBA image to disk.
	outFile, err := os.Create("out.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	bio := bufio.NewWriter(outFile)
	err = png.Encode(bio, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = bio.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")

	return string("FOI..."), nil
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
