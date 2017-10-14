package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

func Thresholding(img *image.Gray, threshold uint8) (out *image.Gray) {
	// get image bounds
	rectangle := (*img).Bounds()

	// create output image
	out = image.NewGray(rectangle)

	for y := rectangle.Min.Y; y < rectangle.Max.Y; y++ {
		for x := rectangle.Min.X; x < rectangle.Max.X; x++ {

			greyValue := (*img).GrayAt(x, y)
			if greyValue.Y <= threshold {
				greyValue.Y = 0
			} else {
				greyValue.Y = 255
			}

			out.SetGray(x, y, greyValue)
		}
	}
	return
}

func RandomDithering(img *image.Gray) (out *image.Gray) {
	// get image bounds
	rectangle := (*img).Bounds()

	// create output image
	out = image.NewGray(rectangle)

	rand.Seed(time.Now().Unix())

	for y := rectangle.Min.Y; y < rectangle.Max.Y; y++ {
		for x := rectangle.Min.X; x < rectangle.Max.X; x++ {
			var threshold uint8 = uint8(rand.Intn(256))

			greyValue := (*img).GrayAt(x, y)
			if greyValue.Y <= threshold {
				greyValue.Y = 0
			} else {
				greyValue.Y = 255
			}

			out.SetGray(x, y, greyValue)
		}
	}
	return
}

func OrderedDithering(img *image.Gray, matrix [][]float64) (out *image.Gray) {

	n := len(matrix)

	// get image bounds
	rectangle := (*img).Bounds()

	// create output image
	out = image.NewGray(rectangle)

	for y := rectangle.Min.Y; y < rectangle.Max.Y; y++ {
		for x := rectangle.Min.X; x < rectangle.Max.X; x++ {

			greyValue := (*img).GrayAt(x, y)
			// to [0, 1]
			var value float64 = float64(greyValue.Y) / 255.0

			// nearest color and other color
			firstColor := 1.0
			secondColor := 0.0

			if value < 0.5 {
				firstColor = 0.0
				secondColor = 1.0
			}

			matrixValue := matrix[x%n][y%n]

			var resultValue float64

			if distance := math.Abs(firstColor - value); distance < matrixValue {
				resultValue = firstColor
			} else {
				resultValue = secondColor
			}

			greyValue.Y = uint8(math.Floor(255 * resultValue))
			out.SetGray(x, y, greyValue)
		}
	}
	return
}

func LineErrorDiffusionDithering(img *image.Gray, threshold uint8) (out *image.Gray) {
	// get image bounds
	rectangle := (*img).Bounds()

	// create output image
	out = image.NewGray(rectangle)

	rand.Seed(time.Now().Unix())

	var error uint8
	for y := rectangle.Min.Y; y < rectangle.Max.Y; y++ {
		error = 0
		for x := rectangle.Min.X; x < rectangle.Max.X; x++ {

			greyValue := (*img).GrayAt(x, y)

			value := greyValue.Y

			if value <= threshold {
				error = error + value
				greyValue.Y = 0
			} else {
				error = error + value - 255
				greyValue.Y = 255
			}

			out.SetGray(x, y, greyValue)
		}
	}
	return
}

func MakeOrderedMatrix() (matrix [][]float64) {

	matrix = [][]float64{
		{0, 8, 2, 10},
		{12, 4, 14, 6},
		{3, 11, 1, 9},
		{15, 7, 13, 5}}

	n := 4
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			matrix[i][j] = matrix[i][j] / 16
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

			fullColor := (*img).At(x, y)
			c := color.GrayModel.Convert(fullColor).(color.Gray)
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

	// output_img = Thresholding(output_img, 100)

	// output_img = RandomDithering(output_img)

	// output_img = OrderedDithering(output_img, MakeOrderedMatrix())

	output_img = LineErrorDiffusionDithering(output_img, 100)

	png.Encode(output_img_raw, image.Image(output_img))
}
