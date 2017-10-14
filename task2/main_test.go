package main

import (
	"image/png"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	// load file
	input_img_raw, _ := os.Open("input.png")

	// decode png image
	input_img, _ := png.Decode(input_img_raw)

	output_img := MakeImageGray(&input_img)
	output_img = ApproximateNearest(output_img)
}
