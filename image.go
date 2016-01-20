package main

import (
	"encoding/hex"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math/rand"
)

var font *truetype.Font
var cjkFont *truetype.Font
var colors = []color.RGBA{
	{191, 210, 215, 255},
	{87, 130, 139, 255},
	{87, 130, 139, 255},
	{29, 137, 160, 255},
	{31, 112, 129, 255},
	{141, 166, 174, 255},
	{110, 163, 175, 255},
	{15, 85, 100, 255},
}

func initFont() {
	fontBytes, err := ioutil.ReadFile("./font.ttf")
	if err != nil {
		fmt.Println(err)
	}

	cjkFontBytes, err := ioutil.ReadFile("./wqy-zenhei.ttc")
	if err != nil {
		fmt.Println(err)
	}

	font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println(err)
	}

	cjkFont, err = freetype.ParseFont(cjkFontBytes)
	if err != nil {
		fmt.Println(err)
	}
}

func hexToRGB(hexString string) color.RGBA {
	// TODO: strip #, allow 3-value hex colors
	values, err := hex.DecodeString(hexString)
	if err != nil {
		fmt.Println(err)
	}
	return color.RGBA{values[0], values[1], values[2], 255}
}

func generateImage(params Params) image.Image {
	var detectedFont *truetype.Font

	if []rune(params.text)[0] > '\u2E7F' {
		detectedFont = cjkFont
	} else {
		detectedFont = font
	}

	fontSize := float64(params.size / 2)

	img := image.NewRGBA(image.Rect(0, 0, params.size, params.size))

	if params.seed != 0 {
		rand.Seed(params.seed)
	}

	var bgColor color.RGBA
	if params.color == (color.RGBA{}) {
		bgColor = colors[rand.Intn(len(colors))]
	} else {
		bgColor = params.color
	}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(detectedFont)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.White)

	fontWidth := (font.Bounds(c.PointToFixed(fontSize) >> 6)).Max.X

	i := int(float64(params.size) - (0.75 * float64(fontWidth) * float64(len(params.text))))
	pt := freetype.Pt(i/2, int(fontSize*1.333))
	_, err := c.DrawString(params.text, pt)
	if err != nil {
		fmt.Println(err)
	}

	return img
}
