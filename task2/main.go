package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {

	// load file
	old_img_raw, _ := os.Open("img.png")

	// decode png image
	old_img, _ := png.Decode(old_img_raw)

	// get image bounds
	rectangle := old_img.Bounds()

	// create new file
	new_img_raw, _ := os.Create("new.png")

	// create new image
	new_img := image.NewGray(rectangle)

	for y := rectangle.Min.Y; y < rectangle.Max.Y; y++ {
		for x := rectangle.Min.X; x < rectangle.Max.X; x++ {

			c := color.GrayModel.Convert(old_img.At(x, y)).(color.Gray)
			new_img.Set(x, y, c)
		}
	}
	png.Encode(new_img_raw, new_img)
}
