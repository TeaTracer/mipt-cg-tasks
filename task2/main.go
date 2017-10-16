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

	var err uint8
	var moveRight, errIsPositive bool
	var nX, nY int
	nX = rectangle.Max.X
	nY = rectangle.Max.Y

	for y := rectangle.Min.Y; y < nY; y++ {

		if y%2 == 0 {
			moveRight = true
		} else {
			moveRight = false
		}

		for x := rectangle.Min.X; x < nX; x++ {

			greyValue := (*img).GrayAt(x, y)
			value := greyValue.Y

			if value <= threshold {
				greyValue.Y = 0
				err = value
				errIsPositive = true
				applyFSError(img, x, y, nX-1, nY-1, err, moveRight, errIsPositive)
			} else {
				greyValue.Y = 255
				err = 255 - value
				errIsPositive = false
				applyFSError(img, x, y, nX-1, nY-1, err, moveRight, errIsPositive)
			}

			out.SetGray(x, y, greyValue)
		}
	}
	return
}

func cropValue(value, err uint8, errIsPositive bool) (new_value uint8) {
	var new_value_int int

	if errIsPositive == true {
		new_value_int = int(value + err)
		if new_value_int > 255 {
			new_value = 255
		} else {
			new_value = uint8(new_value_int)
		}
	} else {
		new_value_int = int(value - err)
		if new_value_int < 0 {
			new_value = 0
		} else {
			new_value = uint8(new_value_int)
		}
	}
	return
}

func applyFSError(img *image.Gray, x, y, maxX, maxY int, err uint8, moveRight, errIsPositive bool) {
	if x == 0 || x == maxX || y == maxY {
		return
	}

	var value, applied_err, new_value uint8
	var greyValue color.Gray
	var coords [4][2]int
	var coefs [4]float64
	var coef float64
	var xx, yy int

	if moveRight == true {
		coords = [4][2]int{{x + 1, y}, {x + 1, y + 1}, {x, y + 1}, {x - 1, y + 1}}
	} else {
		coords = [4][2]int{{x - 1, y}, {x - 1, y + 1}, {x, y + 1}, {x + 1, y + 1}}
	}

	coefs = [4]float64{7.0 / 16.0, 1.0 / 16.0, 5.0 / 16.0, 3.0 / 16.0}

	for i, coord := range coords {
		xx = coord[0]
		yy = coord[1]
		coef = coefs[i]

		greyValue = (*img).GrayAt(xx, yy)
		value = greyValue.Y
		applied_err = uint8(math.Floor(float64(err) * coef))
		new_value = cropValue(value, applied_err, errIsPositive)
		greyValue.Y = new_value
		(*img).SetGray(xx, yy, greyValue)

	}
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
		output_img = FloydSteinbergDithering(output_img, 100)
	default:
		output_img = Thresholding(output_img, 100)
	}

	png.Encode(output_img_raw, image.Image(output_img))
}
