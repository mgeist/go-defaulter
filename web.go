package main

import (
	"fmt"
	"html/template"
	"image/png"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var templates = template.Must(template.ParseFiles("test.html"))

func parseParams(params url.Values) (int, int64, string) {
	sizeString := params.Get("size")
	seedString := params.Get("seed")
	text := params.Get("text")

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

	return size, int64(seed), text
}

func getPort() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return port
}

func handler(w http.ResponseWriter, r *http.Request) {
	size, seed, text := parseParams(r.URL.Query())
	png.Encode(w, generateImage(size, seed, text))
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
