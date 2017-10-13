package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	input_img_ptr := flag.String("i", "input.png", "Input PNG image.")
	output_img_ptr := flag.String("o", "output.png", "Output PNG image.")
	flag.Parse()

	// load file
	input_img_raw, _ := os.Open(*input_img_ptr)

	// decode png image
	input_img, _ := png.Decode(input_img_raw)

	// get image bounds
	rectangle := input_img.Bounds()

	// create output file
	output_img_raw, _ := os.Create(*output_img_ptr)

	// create output image
	output_img := image.NewGray(rectangle)

	for y := rectangle.Min.Y; y < rectangle.Max.Y; y++ {
		for x := rectangle.Min.X; x < rectangle.Max.X; x++ {

			c := color.GrayModel.Convert(input_img.At(x, y)).(color.Gray)
			output_img.Set(x, y, c)
		}
	}
	png.Encode(output_img_raw, output_img)
}
