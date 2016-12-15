package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/", Handler)

	http.ListenAndServe(":8080", nil)
}
func Handler(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	body, readErr := ioutil.ReadAll(req.Body)
	if readErr != nil {
		fmt.Fprintf(res, "error reading: %s!", readErr)
		return
	}

	indexHead := bytes.Index(body, []byte("\r\n\r\n")) + 4
	bodyDecode := bytes.NewReader(body[indexHead:])
	img, imageErr := png.Decode(bodyDecode)
	if imageErr != nil {
		fmt.Fprintf(res, "error image: %s!%s", imageErr, body[indexHead:])
		return
	}
	superImage := superSampling(img)

	res.Header().Set("Content-Type", "image/png")
	png.Encode(res, superImage)
}

func superSampling(source image.Image) image.Image {
	img := image.NewRGBA(source.Bounds())
	for py := 0; py < img.Rect.Dy(); py++ {
		for px := 0; px < img.Rect.Dx(); px++ {
			img.Set(px, py, average(around(source, px, py)))
		}
	}
	return img
}
func around(source image.Image, px, py int) [9]color.Color {
	var colors [9]color.Color
	colors[0] = source.At(px-1, py-1)
	colors[1] = source.At(px, py-1)
	colors[2] = source.At(px+1, py-1)
	colors[3] = source.At(px-1, py)
	colors[4] = source.At(px, py)
	colors[5] = source.At(px+1, py)
	colors[6] = source.At(px-1, py+1)
	colors[7] = source.At(px, py+1)
	colors[8] = source.At(px+1, py+1)
	return colors
}

func average(colors [9]color.Color) color.Color {
	var r, g, b, n uint32
	for _, c := range colors {
		cr, cg, cb, _ := c.RGBA()
		if cr == 0 && cg == 0 && cb == 0 {
			continue
		}
		r += cr
		g += cg
		b += cb
		n++
	}
	if n == 0 {
		return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
	}
	return color.RGBA{uint8(r / n >> 8), uint8(g / n >> 8), uint8(b / n >> 8), 255}
}
