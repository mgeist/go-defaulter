package main

import (
	"fmt"
	"image/png"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

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

	size = int(math.Max(math.Min(2048, float64(size)), 25))

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

func main() {
	initFont()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+getPort(), nil)
}
