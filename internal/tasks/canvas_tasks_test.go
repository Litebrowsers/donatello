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
	"testing"
)

var CanvasWidth = 20
var CanvasHeight = 20
var CanvasSize = 20

func TestRectangle_Encode(t *testing.T) {
	rec := Rectangle{Color: "FF0000", W: 5, H: 3, X: 10, Y: 5}
	expected := "R:FF0000:5:3:10:5"
	if rec.Encode() != expected {
		t.Errorf("Rectangle.Encode() failed. Expected %s, got %s", expected, rec.Encode())
	}
}

func TestCircle_Encode(t *testing.T) {
	circ := Circle{Color: "00FF00", R: 4, X: 15, Y: 15}
	expected := "C:00FF00:4:15:15"
	if circ.Encode() != expected {
		t.Errorf("Circle.Encode() failed. Expected %s, got %s", expected, circ.Encode())
	}
}

func TestTriangle_Encode(t *testing.T) {
	tri := Triangle{Color: "0000FF", X1: 2, Y1: 2, X2: 6, Y2: 2, X3: 4, Y3: 6}
	expected := "T:0000FF:2:2:6:2:4:6"
	if tri.Encode() != expected {
		t.Errorf("Triangle.Encode() failed. Expected %s, got %s", expected, tri.Encode())
	}
}

func TestLine_Encode(t *testing.T) {
	line := Line{Color: "FF00FF", X1: 5, Y1: 5, X2: 12, Y2: 8, Thickness: 2}
	expected := "L:FF00FF:5:5:12:8:2"
	if line.Encode() != expected {
		t.Errorf("Line.Encode() failed. Expected %s, got %s", expected, line.Encode())
	}
}

func TestEllipse_Encode(t *testing.T) {
	ellipse := Ellipse{Color: "FFFF00", RX: 6, RY: 3, X: 20, Y: 10}
	expected := "E:FFFF00:6:3:20:10"
	if ellipse.Encode() != expected {
		t.Errorf("Ellipse.Encode() failed. Expected %s, got %s", expected, ellipse.Encode())
	}
}

func TestTaskGenerator_GenerateTask(t *testing.T) {
	rec := Rectangle{Color: "FF0000", W: 5, H: 3, X: 10, Y: 5}
	circ := Circle{Color: "00FF00", R: 4, X: 15, Y: 15}
	tri := Triangle{Color: "0000FF", X1: 2, Y1: 2, X2: 6, Y2: 2, X3: 4, Y3: 6}

	tg := NewTaskGenerator(rec, circ, tri)
	expected := "R:FF0000:5:3:10:5;C:00FF00:4:15:15;T:0000FF:2:2:6:2:4:6"

	if tg.GenerateTask() != expected {
		t.Errorf("TaskGenerator.GenerateTask() failed. Expected %s, got %s", expected, tg.GenerateTask())
	}

	// Test with AddShape
	line := Line{Color: "FF00FF", X1: 5, Y1: 5, X2: 12, Y2: 8, Thickness: 2}
	tg.AddShape(line)
	expectedWithLine := "R:FF0000:5:3:10:5;C:00FF00:4:15:15;T:0000FF:2:2:6:2:4:6;L:FF00FF:5:5:12:8:2"
	if tg.GenerateTask() != expectedWithLine {
		t.Errorf("TaskGenerator.GenerateTask() with AddShape failed. Expected %s, got %s", expectedWithLine, tg.GenerateTask())
	}
}

func TestGenerateRandomShapes(t *testing.T) {
	count := 5
	shapes := GenerateRandomShapes(CanvasSize, count)

	if len(shapes) != count {
		t.Errorf("GenerateRandomShapes() returned %d shapes, expected %d", len(shapes), count)
	}

	for i, s := range shapes {
		if s.Encode() == "" {
			t.Errorf("Shape at index %d returned empty encoded string", i)
		}
	}

	// Test with count = 0
	shapesZero := GenerateRandomShapes(CanvasSize, 0)
	if shapesZero != nil {
		t.Errorf("GenerateRandomShapes(0) should return nil, got %v", shapesZero)
	}
}

func TestGenerateRandomEvenSizedPrimitives(t *testing.T) {
	count := 10 // Increase count to have a higher chance of generating both types
	primitives := GenerateRandomEvenSizedPrimitives(CanvasSize, count)

	if len(primitives) != count {
		t.Errorf("GenerateRandomEvenSizedPrimitives() returned %d primitives, expected %d", len(primitives), count)
	}

	for i, s := range primitives {
		if s.Encode() == "" {
			t.Errorf("Primitive at index %d returned empty encoded string", i)
		}

		switch v := s.(type) {
		case Rectangle:
			if v.W != v.H {
				t.Errorf("Square at index %d has unequal width and height: W=%d, H=%d", i, v.W, v.H)
			}
			if v.W%2 != 0 {
				t.Errorf("Square at index %d has odd side length: %d", i, v.W)
			}
		case Line:
			if v.X1 != v.X2 && v.Y1 != v.Y2 {
				t.Errorf("Line at index %d is neither vertical nor horizontal: (%d,%d)-(%d,%d)", i, v.X1, v.Y1, v.X2, v.Y2)
			}
		default:
			t.Errorf("Unknown primitive type at index %d: %T", i, s)
		}
	}

	// Test with count = 0
	primitivesZero := GenerateRandomEvenSizedPrimitives(CanvasSize, 0)
	if primitivesZero != nil {
		t.Errorf("GenerateRandomEvenSizedPrimitives(0) should return nil, got %v", primitivesZero)
	}
}

func TestMatrixAndHashGeneration(t *testing.T) {
	// 1. Create a canvas
	canvas := NewCanvas(CanvasWidth, CanvasHeight)

	// 2. Define and draw a shape
	// Using a rectangle similar to the example from gemini.md
	// Color #6F79D2 -> R=111, G=121, B=210
	rect := Rectangle{
		Color: "6F79D2",
		X:     5,
		Y:     0,
		W:     4,
		H:     6,
	}
	err := canvas.DrawShapes([]Shape{rect})
	if err != nil {
		return
	}

	// 3. Create the expected matrix for the Red channel (value 111)
	expectedRedPix := make([]uint8, CanvasWidth*CanvasHeight)
	for y := 0; y < 6; y++ {
		for x := 5; x < 9; x++ {
			expectedRedPix[y*CanvasWidth+x] = 111
		}
	}

	// 4. Calculate the expected hash for the Red channel
	expectedRedHash, err := calculateHash(expectedRedPix)
	if err != nil {
		t.Fatalf("Failed to calculate expected red hash: %v", err)
	}

	// 5. Get the actual hashes from the canvas
	actualHashes, err := canvas.CalculateHashes()
	if err != nil {
		t.Fatalf("Canvas.CalculateHashes() returned an error: %v", err)
	}

	// 6. Compare the hashes
	if actualHashes["red"] != expectedRedHash {
		t.Errorf("Red channel hash mismatch. Expected %s, got %s", expectedRedHash, actualHashes["red"])
	}

	// Optional: You can do the same for other channels (e.g., Green)
	expectedGreenPix := make([]uint8, CanvasWidth*CanvasHeight)
	for y := 0; y < 6; y++ {
		for x := 5; x < 9; x++ {
			expectedGreenPix[y*CanvasWidth+x] = 121
		}
	}
	expectedGreenHash, err := calculateHash(expectedGreenPix)
	if err != nil {
		t.Fatalf("Failed to calculate expected green hash: %v", err)
	}
	if actualHashes["green"] != expectedGreenHash {
		t.Errorf("Green channel hash mismatch. Expected %s, got %s", expectedGreenHash, actualHashes["green"])
	}
}

func TestSpecificLineShapeHashes(t *testing.T) {
	shapes := []Shape{
		Line{Color: "CE818F", X1: 13, Y1: 10, X2: 13, Y2: 13, Thickness: 4},
	}

	canvas := NewCanvas(CanvasWidth, CanvasHeight)
	err := canvas.DrawShapes(shapes)
	if err != nil {
		return
	}

	actualHashes, err := canvas.CalculateHashes()
	if err != nil {
		t.Fatalf("Canvas.CalculateHashes() returned an error: %v", err)
	}

	expectedHashes := map[string]string{
		"red":   "d144f3a83b769ff5ec0399e60c3bfd9d6e21af518e3a7817350f40ae26d4d077",
		"green": "a2bc793bb323f2ab2a91421d17f96e20e6e3c2c029a9fc8c64b2edd8c8c47a4e",
		"blue":  "8da048303ea3f18e92a3776baf1ad15c424e200477114b61f7a9e0c3c17943a1",
		"alpha": "45f282d79130a1ccc85813697a7c52df38d9fbf6d43ceeaba5553821d321063b",
	}

	if actualHashes["red"] != expectedHashes["red"] {
		t.Errorf("Red channel hash mismatch.\nExpected: %s\nActual:   %s", expectedHashes["red"], actualHashes["red"])
	}
	if actualHashes["green"] != expectedHashes["green"] {
		t.Errorf("Green channel hash mismatch.\nExpected: %s\nActual:   %s", expectedHashes["green"], actualHashes["green"])
	}
	if actualHashes["blue"] != expectedHashes["blue"] {
		t.Errorf("Blue channel hash mismatch.\nExpected: %s\nActual:   %s", expectedHashes["blue"], actualHashes["blue"])
	}
	if actualHashes["alpha"] != expectedHashes["alpha"] {
		t.Errorf("Alpha channel hash mismatch.\nExpected: %s\nActual:   %s", expectedHashes["alpha"], actualHashes["alpha"])
	}
}
