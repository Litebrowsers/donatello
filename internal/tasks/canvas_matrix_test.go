/*
# Donatello

Copyright © 2025 Litebrowsers
Licensed under a Proprietary License

This software is the confidential and proprietary information of Litebrowsers
Unauthorized copying, redistribution, or use is prohibited.
For licensing inquiries, contact:
vera cohopie at gmail dot com
thor betson at gmail dot com
*/

package tasks

import (
	"image/color"
	"testing"
)

func TestNewCanvas(t *testing.T) {
	width, height := 100, 200
	canvas := NewCanvas(width, height)

	if canvas.R.Bounds().Dx() != width || canvas.R.Bounds().Dy() != height {
		t.Errorf("NewCanvas() created a canvas with incorrect dimensions for R channel")
	}
	if canvas.G.Bounds().Dx() != width || canvas.G.Bounds().Dy() != height {
		t.Errorf("NewCanvas() created a canvas with incorrect dimensions for G channel")
	}
	if canvas.B.Bounds().Dx() != width || canvas.B.Bounds().Dy() != height {
		t.Errorf("NewCanvas() created a canvas with incorrect dimensions for B channel")
	}
	if canvas.A.Bounds().Dx() != width || canvas.A.Bounds().Dy() != height {
		t.Errorf("NewCanvas() created a canvas with incorrect dimensions for A channel")
	}
}

func TestCanvas_DrawShapes_Rectangle(t *testing.T) {
	canvas := NewCanvas(20, 20)
	rect := Rectangle{Color: "FF0000", W: 5, H: 3, X: 10, Y: 5}
	err := canvas.DrawShapes([]Shape{rect})
	if err != nil {
		return
	}

	// Check a pixel inside the rectangle
	c := canvas.R.At(12, 6).(color.Gray)
	if c.Y != 255 {
		t.Errorf("Expected pixel to be red (255), but got %d", c.Y)
	}

	// Check a pixel outside the rectangle
	c = canvas.R.At(0, 0).(color.Gray)
	if c.Y != 0 {
		t.Errorf("Expected pixel to be black (0), but got %d", c.Y)
	}
}

func TestCanvas_DrawShapes_Line(t *testing.T) {
	canvas := NewCanvas(20, 20)
	line := Line{Color: "00FF00", X1: 5, Y1: 5, X2: 15, Y2: 5, Thickness: 2}
	err := canvas.DrawShapes([]Shape{line})
	if err != nil {
		return
	}

	// Check a pixel on the line
	c := canvas.G.At(10, 5).(color.Gray)
	if c.Y != 255 {
		t.Errorf("Expected pixel to be green (255), but got %d", c.Y)
	}

	// Check a pixel off the line
	c = canvas.G.At(0, 0).(color.Gray)
	if c.Y != 0 {
		t.Errorf("Expected pixel to be black (0), but got %d", c.Y)
	}
}

func Test_hexToRGBA(t *testing.T) {
	rgba, _ := hexToRGBA("FF0000")
	if rgba.R != 255 || rgba.G != 0 || rgba.B != 0 || rgba.A != 255 {
		t.Errorf("hexToRGBA('FF0000') failed. Expected (255, 0, 0, 255), got (%d, %d, %d, %d)", rgba.R, rgba.G, rgba.B, rgba.A)
	}

	rgba, _ = hexToRGBA("00FF00")
	if rgba.R != 0 || rgba.G != 255 || rgba.B != 0 || rgba.A != 255 {
		t.Errorf("hexToRGBA('00FF00') failed. Expected (0, 255, 0, 255), got (%d, %d, %d, %d)", rgba.R, rgba.G, rgba.B, rgba.A)
	}

	rgba, _ = hexToRGBA("0000FF")
	if rgba.R != 0 || rgba.G != 0 || rgba.B != 255 || rgba.A != 255 {
		t.Errorf("hexToRGBA('0000FF') failed. Expected (0, 0, 255, 255), got (%d, %d, %d, %d)", rgba.R, rgba.G, rgba.B, rgba.A)
	}
}
