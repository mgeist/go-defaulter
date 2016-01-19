package main

import (
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

	font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println(err)
	}
}

func generateImage(size int, seed int64, text string) image.Image {
	fontSize := float64(size / 2)

	img := image.NewRGBA(image.Rect(0, 0, size, size))

	if seed != 0 {
		rand.Seed(seed)
	}
	randColor := colors[rand.Intn(len(colors))]
	draw.Draw(img, img.Bounds(), &image.Uniform{randColor}, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.White)

	fontWidth := (font.Bounds(c.PointToFixed(fontSize) >> 6)).Max.X

	i := int(float64(size) - (0.75 * float64(fontWidth) * float64(len(text))))
	pt := freetype.Pt(i/2, int(fontSize*1.333))
	_, err := c.DrawString(text, pt)
	if err != nil {
		fmt.Println(err)
	}

	return img
}
