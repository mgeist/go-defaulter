package main

import (
	"encoding/hex"
	"fmt"
	"github.com/mgeist/draw2d/draw2dimg"
	"github.com/mgeist/freetype"
	"github.com/mgeist/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"
	"math/rand"
)

var defaultFont *truetype.Font
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

	defaultFont, err = freetype.ParseFont(fontBytes)
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
	fgColor := color.RGBA{255, 255, 255, 255}

	var detectedFont *truetype.Font
	if []rune(params.text)[0] > '\u2E7F' {
		detectedFont = cjkFont
	} else {
		detectedFont = defaultFont
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

	if params.border {
		bgColor, fgColor = fgColor, bgColor
	}

	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.ZP, draw.Src)

	if params.border {
		strokeWidth := fontSize * 0.08
		circleSize := fontSize * 0.92
		arcAngle := math.Pi * 2

		gc := draw2dimg.NewGraphicContext(img)
		gc.SetStrokeColor(fgColor)
		gc.SetLineWidth(strokeWidth)
		gc.ArcTo(fontSize, fontSize, circleSize, circleSize, arcAngle, arcAngle)
		gc.Stroke()
	}

	d := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(fgColor),
		Face: truetype.NewFace(detectedFont, &truetype.Options{
			Size: fontSize,
			DPI:  72,
		}),
	}
	d.Dot = fixed.Point26_6{
		X: (fixed.I(params.size) - d.MeasureString(params.text)) / 2,
		Y: fixed.I(int(math.Ceil(fontSize * 1.35))),
	}
	d.DrawString(params.text)

	return img
}

func generatePie(params PieParams) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, params.size, params.size))

	progressColor := progressColors[params.color]
	halfWidth := float64(params.size / 2)
	strokeWidth := math.Floor(halfWidth * 0.1)
	circleSize := halfWidth * 0.94
	arcAngle := math.Pi * 2

	gc := draw2dimg.NewGraphicContext(img)
	gc.BeginPath()
	gc.SetStrokeColor(progressColor)
	gc.SetFillColor(color.RGBA{255, 255, 255, 255})
	gc.SetLineWidth(strokeWidth)
	gc.ArcTo(halfWidth, halfWidth, circleSize, circleSize, arcAngle, arcAngle)
	gc.Close()
	gc.FillStroke()

	startAngle := 270 * (math.Pi / 180.0)
	angle := (360 * (float64(params.progress) * 0.01)) * (math.Pi / 180.0)
	gc.SetFillColor(progressColor)
	gc.BeginPath()
	gc.MoveTo(halfWidth, halfWidth)
	gc.ArcTo(halfWidth, halfWidth, circleSize, circleSize, startAngle, angle)
	gc.Close()
	gc.Fill()

	return img
}

func generateHorseshoe(params PieParams) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, params.size, params.size))

	gray := color.RGBA{247, 246, 245, 255}
	progressColor := progressColors[params.color]
	progressShadow := progressShadows[params.color]

	halfWidth := float64(params.size / 2)
	strokeWidth := float64(params.size) * 0.086
	bigStrokeWidth := float64(params.size) * 0.1
	circleSize := halfWidth * 0.71
	bigCircleSize := float64(params.size) * 0.4
	startAngle := 120 * (math.Pi / 180.0)
	angle := (300 * (float64(100) * 0.01)) * (math.Pi / 180.0)

	gc := draw2dimg.NewGraphicContext(img)
	gc.BeginPath()
	gc.SetStrokeColor(gray)
	gc.SetLineWidth(strokeWidth)
	gc.ArcTo(halfWidth, halfWidth, circleSize, circleSize, startAngle, angle)
	gc.Stroke()

	angle = (300 * (float64(params.progress) * 0.01)) * (math.Pi / 180.0)
	gc.BeginPath()
	gc.SetStrokeColor(progressShadow)
	gc.SetLineWidth(strokeWidth)
	gc.ArcTo(halfWidth, halfWidth, circleSize, circleSize, startAngle, angle)
	gc.Stroke()

	gc.BeginPath()
	gc.SetStrokeColor(progressColor)
	gc.SetLineWidth(bigStrokeWidth)
	gc.ArcTo(halfWidth, halfWidth, bigCircleSize, bigCircleSize, startAngle, angle)
	gc.Stroke()

	return img
}
