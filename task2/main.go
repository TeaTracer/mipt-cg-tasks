package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

func ApproximateThresholding(img *image.Gray, threshold uint8) (out *image.Gray) {
	// get image bounds
	rectangle := (*img).Bounds()

	// create output image
	out = image.NewGray(rectangle)

	for y := rectangle.Min.Y; y < rectangle.Max.Y; y++ {
		for x := rectangle.Min.X; x < rectangle.Max.X; x++ {

			grey_value := (*img).GrayAt(x, y)
			if grey_value.Y <= threshold {
				grey_value.Y = 0
			} else {
				grey_value.Y = 255
			}

			out.SetGray(x, y, grey_value)
		}
	}
	return
}

func ApproximateRandomDithering(img *image.Gray) (out *image.Gray) {
	// get image bounds
	rectangle := (*img).Bounds()

	// create output image
	out = image.NewGray(rectangle)

	rand.Seed(time.Now().Unix())

	for y := rectangle.Min.Y; y < rectangle.Max.Y; y++ {
		for x := rectangle.Min.X; x < rectangle.Max.X; x++ {
			var threshold uint8 = uint8(rand.Intn(256))

			grey_value := (*img).GrayAt(x, y)
			if grey_value.Y <= threshold {
				grey_value.Y = 0
			} else {
				grey_value.Y = 255
			}

			out.SetGray(x, y, grey_value)
		}
	}
	return
}

func MakeImageGray(img *image.Image) (out *image.Gray) {

	// get image bounds
	rectangle := (*img).Bounds()

	// create output image
	out = image.NewGray(rectangle)

	// iterate over all points
	for y := rectangle.Min.Y; y < rectangle.Max.Y; y++ {
		for x := rectangle.Min.X; x < rectangle.Max.X; x++ {

			full_color := (*img).At(x, y)
			c := color.GrayModel.Convert(full_color).(color.Gray)
			out.Set(x, y, c)
		}
	}
	return
}

func main() {
	input_img_ptr := flag.String("i", "input.png", "Input PNG image.")
	output_img_ptr := flag.String("o", "output.png", "Output PNG image.")
	// approximation_type := flag.String("a", "default", "Approximation type (default='nearest')")
	flag.Parse()

	// load file
	input_img_raw, _ := os.Open(*input_img_ptr)

	// decode png image
	input_img, _ := png.Decode(input_img_raw)

	// create output file
	output_img_raw, _ := os.Create(*output_img_ptr)

	output_img := MakeImageGray(&input_img)
	// output_img = ApproximateThresholding(output_img, 100)
	output_img = ApproximateRandomDithering(output_img)
	png.Encode(output_img_raw, image.Image(output_img))
}
