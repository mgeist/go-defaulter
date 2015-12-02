package main

import (
  "image"
  "image/color"
  "image/jpeg"
  "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
  img := image.NewRGBA(image.Rect(0,0,100,100))

	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			img.Set(x, y, color.RGBA{0x88,0xff,0x88,0xff})
		}
	}

	jpeg.Encode(w, img, &jpeg.Options{80})
}

func main() {
  http.HandleFunc("/", handler)
  http.ListenAndServe(":8080", nil)
}
