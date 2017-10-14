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

			matrixValue := matrix[x%n][y%n]

			var resultValue float64

			if value <= matrixValue {
				resultValue = 0.0
			} else {
				resultValue = 1.0
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

	var err int
	for y := rectangle.Min.Y; y < rectangle.Max.Y; y++ {
		err = 0
		for x := rectangle.Min.X; x < rectangle.Max.X; x++ {

			greyValue := (*img).GrayAt(x, y)

			value := int(greyValue.Y)

			if value+err <= int(threshold) {
				err = err + value
				greyValue.Y = 0
			} else {
				err = err + value - 255
				greyValue.Y = 255
			}

			out.SetGray(x, y, greyValue)
		}
	}
	return
}

func LineAlternationErrorDiffusionDithering(img *image.Gray, threshold uint8) (out *image.Gray) {
	// get image bounds
	rectangle := (*img).Bounds()

	// create output image
	out = image.NewGray(rectangle)

	rand.Seed(time.Now().Unix())

	var err int
	for y := rectangle.Min.Y; y < rectangle.Max.Y; y++ {
		err = 0
		var new_x int
		for x := rectangle.Min.X; x < rectangle.Max.X; x++ {
			if y%2 == 1 {
				new_x = rectangle.Max.X - (x - rectangle.Min.X)
			} else {
				new_x = x
			}

			greyValue := (*img).GrayAt(new_x, y)
			value := int(greyValue.Y)

			if value+err <= int(threshold) {
				err = err + value
				greyValue.Y = 0
			} else {
				err = err + value - 255
				greyValue.Y = 255
			}

			out.SetGray(new_x, y, greyValue)
		}
	}
	return
}

func FloydSteinbergDithering(img *image.Gray, threshold uint8) (out *image.Gray) {
	// get image bounds
	rectangle := (*img).Bounds()

	// create output image
	out = image.NewGray(rectangle)

	rand.Seed(time.Now().Unix())

	var err int
	var nX, nY int
	nX = rectangle.Max.X
	nY = rectangle.Max.Y

	for y := rectangle.Min.Y; y < nY; y++ {
		err = 0
		for x := rectangle.Min.X; x < nX; x++ {

			greyValue := (*img).GrayAt(x, y)
			value := int(getFloydSteinbergValue(img, x, y, nX-1, nY-1))

			if value+err <= int(threshold) {
				err = err + value
				greyValue.Y = 0
			} else {
				err = err + value - 255
				greyValue.Y = 255
			}

			out.SetGray(x, y, greyValue)
		}
	}
	return
}

func getFSValue(ax, ay, by, cy uint8) (val uint8) {
	var axi, ayi, byi, cyi int
	axi = int(ax)
	ayi = int(ay)
	byi = int(by)
	cyi = int(cy)

	val2 := axi*7 + ayi*1 + byi*5 + cyi*3
	val = uint8(math.Floor(float64(val2) / 16.0))
	return
}

func getFloydSteinbergValue(img *image.Gray, x, y, maxX, maxY int) (value uint8) {
	var ax, ay, by, cy uint8
	var axx, axy, ayx, ayy, byx, byy, cyx, cyy int

	axx = x + 1
	axy = y
	ayx = x + 1
	ayy = y - 1
	byx = x
	byy = y - 1
	cyx = x - 1
	cyy = y - 1

	if y == maxY {
		ayy = y
		byy = y
		cyy = y
	}

	if x == maxX {
		axx = x
		ayx = x
	}

	if x == 0 {
		cyx = x
	}

	ax = (*img).GrayAt(axx, axy).Y
	ay = (*img).GrayAt(ayx, ayy).Y
	by = (*img).GrayAt(byx, byy).Y
	cy = (*img).GrayAt(cyx, cyy).Y
	rs := getFSValue(ax, ay, by, cy)
	return rs
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
	approximation_type := flag.String("a", "threshold",
		"Approximation type ['threshold', 'random', 'ordered', 'diffusion', 'floyd'] (default='threshold')")
	flag.Parse()

	// load file
	input_img_raw, _ := os.Open(*input_img_ptr)

	// decode png image
	input_img, _ := png.Decode(input_img_raw)

	// create output file
	output_img_raw, _ := os.Create(*output_img_ptr)

	output_img := MakeImageGray(&input_img)

	switch *approximation_type {
	case "threshold":
		output_img = Thresholding(output_img, 100)
	case "random":
		output_img = RandomDithering(output_img)
	case "ordered":
		output_img = OrderedDithering(output_img, MakeOrderedMatrix())
	case "diffusion":
		output_img = LineErrorDiffusionDithering(output_img, 100)
		// output_img = LineAlternationErrorDiffusionDithering(output_img, 100)
	case "floyd":
		output_img = FloydSteinbergDithering(output_img, 150)
	default:
		output_img = Thresholding(output_img, 100)
	}

	png.Encode(output_img_raw, image.Image(output_img))
}
