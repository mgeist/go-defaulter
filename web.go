package main

import (
	"fmt"
	"html/template"
	"image/color"
	"image/png"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var templates = template.Must(template.ParseFiles("test.html"))

var progressColors = map[string]color.RGBA{
	"green":  color.RGBA{130, 187, 65, 255},
	"yellow": color.RGBA{243, 228, 110, 255},
	"red":    color.RGBA{227, 112, 104, 255},
}

type Params struct {
	size   int
	seed   int64
	text   string
	border bool
	color  color.RGBA
}

type PieParams struct {
	size     int
	progress int
	color    color.RGBA
}

func parsePieParams(urlParams url.Values) PieParams {
	sizeString := urlParams.Get("size")
	progressString := urlParams.Get("progress")
	colorString := urlParams.Get("color")

	if len(sizeString) == 0 {
		sizeString = "200"
	}

	size, err := strconv.Atoi(sizeString)
	if err != nil {
		fmt.Println(err)
	}

	size = int(math.Max(math.Min(2048, float64(size)), 1))

	if len(progressString) == 0 {
		progressString = "0"
	}

	progress, err := strconv.Atoi(progressString)
	if err != nil {
		fmt.Println(err)
	}

	progress = int(math.Max(math.Min(100, float64(progress)), 0))

	if len(colorString) == 0 {
		colorString = "green"
	}

	var progressColor color.RGBA
	if c, ok := progressColors[colorString]; ok {
		progressColor = c
	} else {
		progressColor = progressColors["green"]
	}

	params := PieParams{
		size:     size,
		progress: progress,
		color:    progressColor,
	}

	return params
}

func parseParams(urlParams url.Values) Params {
	sizeString := urlParams.Get("size")
	seedString := urlParams.Get("seed")
	text := urlParams.Get("text")
	borderString := urlParams.Get("border")
	hex := urlParams.Get("hex")

	if len(sizeString) == 0 {
		sizeString = "200"
	}

	if len(seedString) == 0 {
		seedString = "0"
	}

	if len(text) == 0 {
		text = "?"
	}

	if len([]rune(text)) > 2 {
		text = text[:2]
	}

	if len(borderString) == 0 {
		borderString = "false"
	}

	size, err := strconv.Atoi(sizeString)
	if err != nil {
		fmt.Println(err)
	}

	size = int(math.Max(math.Min(2048, float64(size)), 1))

	seed, err := strconv.Atoi(seedString)
	if err != nil {
		fmt.Println(err)
	}

	border, err := strconv.ParseBool(borderString)
	if err != nil {
		fmt.Println(err)
	}

	var color color.RGBA
	if len(hex) != 0 {
		color = hexToRGB(hex)
	}

	params := Params{
		size:   size,
		seed:   int64(seed),
		text:   text,
		border: border,
		color:  color,
	}

	return params
}

func getAddr() string {
	port := os.Getenv("PORT")
	var addr string
	if len(port) == 0 {
		addr = "127.0.0.1:8080"
	} else {
		addr = ":" + port
	}
	return addr
}

func handler(w http.ResponseWriter, r *http.Request) {
	params := parseParams(r.URL.Query())
	png.Encode(w, generateImage(params))
}

func pieHandler(w http.ResponseWriter, r *http.Request) {
	params := parsePieParams(r.URL.Query())
	png.Encode(w, generatePie(params))
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	var sizes []int
	var ranges []int

	for i := -10; i <= 110; i += 10 {
		ranges = append(ranges, i)
	}

	for i := -10; i <= 100; i += 20 {
		sizes = append(sizes, i)
	}

	data := struct {
		Sizes  []int
		Ranges []int
	}{
		sizes,
		ranges,
	}

	templates.ExecuteTemplate(w, "test.html", data)
}

func main() {
	initFont()

	http.HandleFunc("/", handler)
	http.HandleFunc("/pie/", pieHandler)
	http.HandleFunc("/test", testHandler)
	// TODO: set proper content headers instead of this
	http.HandleFunc("/favicon.ico", http.NotFound)
	http.ListenAndServe(getAddr(), nil)
}
