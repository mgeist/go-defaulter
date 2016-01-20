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

type Params struct {
	size  int
	seed  int64
	text  string
	color color.RGBA
}

func parseParams(urlParams url.Values) Params {
	sizeString := urlParams.Get("size")
	seedString := urlParams.Get("seed")
	text := urlParams.Get("text")
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

	if len(text) > 2 {
		text = text[:2]
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

	var color color.RGBA
	if len(hex) != 0 {
		color = hexToRGB(hex)
	}

	params := Params{
		size:  size,
		seed:  int64(seed),
		text:  text,
		color: color,
	}

	return params
}

func getPort() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return port
}

func handler(w http.ResponseWriter, r *http.Request) {
	params := parseParams(r.URL.Query())
	png.Encode(w, generateImage(params))
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
	http.HandleFunc("/test", testHandler)
	http.ListenAndServe(":"+getPort(), nil)
}
