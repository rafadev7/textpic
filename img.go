package main

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/image/font"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
)

var (
	dpi      = float64(72.0)
	fontfile = "DejaVuSansMono.ttf"
	hinting  = font.HintingNone // or "none"
	maxSize  = float64(12)      // Font max size
	wonb     = false            // White or Black
	f        *truetype.Font
)

func init() {

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Panicln("Err loading font: ", err.Error())
	}
	f, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.Panicln("Err parsing font: ", err.Error())
	}

}

var grayScale = "#$@B%8&WM*oahkbdpqwmZO0QLCJUYXzcvunxrjft{}[]/\\|()1?-_+~<>i!lI;:,^'''.  "

func ImageToText(url string) ([]string, float64, error) {

	resp, err := http.Get(url)
	if err != nil {
		return []string{}, 0, errors.New("error getting img url:" + err.Error())
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return []string{}, 0, errors.New("error decoding body img:" + err.Error())
	}

	size := float64(img.Bounds().Dy() / 100) // Min of 100 lines
	if size > maxSize {
		size = maxSize
	}

	img = ScaleImage(img, size)

	text := Convert2Ascii(img)

	return text, size, nil
}

func TextToImage(text []string, size float64) ([]byte, error) {
	var err error
	// Initialize the context.
	fg, bg := image.Black, image.White
	if wonb {
		fg, bg = image.White, image.Black
	}

	// How it will advance in the first line we will write
	wx := StringAdvance(text[0], size)
	hx := len(text) * int(size)

	rgba := image.NewRGBA(image.Rect(0, 0, wx, hx)) // 640, 480
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetHinting(hinting)

	// Draw the text.
	pt := freetype.Pt(0, int(size)) // int(c.PointToFixed(>>6)
	for _, line := range text {
		_, err = c.DrawString(line, pt)
		if err != nil {
			return []byte{}, errors.New("error drawing text:" + err.Error())
		}
		pt.Y += c.PointToFixed(size)
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, rgba)
	if err != nil {
		return []byte{}, errors.New("error ecoding img:" + err.Error())
	}

	// Save that RGBA image to disk.
	// outFile, err := os.Create("out.png")
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	// defer outFile.Close()
	// bio := bufio.NewWriter(outFile)
	// err = png.Encode(bio, rgba)
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	// err = bio.Flush()
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Println("Wrote out.png OK.")

	return buf.Bytes(), nil
}

func ScaleImage(img image.Image, size float64) image.Image {
	sz := img.Bounds()
	w := uint(float64(sz.Dx()) / size)
	// Assuming that font's height is about
	// 1.6x highier than font's widht
	h := uint((sz.Dy() * int(w) * 10) / (sz.Dx() * 16))
	img = resize.Resize(w, h, img, resize.Lanczos3)
	return img
}

func Convert2Ascii(img image.Image) []string {
	text := []string{}
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	for i := 0; i < h; i++ {
		line := ""
		for j := 0; j < w; j++ {
			g := color.GrayModel.Convert(img.At(j, i)).(color.Gray)
			pos := int(uint(g.Y) * 70 / 255)
			line = line + string(grayScale[pos])
		}

		text = append(text, line)
	}
	return text
}

func StringAdvance(text string, size float64) int {
	face := truetype.NewFace(f, &truetype.Options{
		Size:    size,
		DPI:     dpi,
		Hinting: hinting,
	})
	return font.MeasureString(face, text).Round()
}
