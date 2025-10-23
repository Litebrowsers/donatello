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

// Package tasks contains the canvas related code.
package tasks

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

// Canvas represents the drawing canvas with separate color channels.
type Canvas struct {
	R, G, B, A *image.Gray
}

// NewCanvas creates a new canvas of the specified width and height.
func NewCanvas(width, height int) *Canvas {
	return &Canvas{
		R: image.NewGray(image.Rect(0, 0, width, height)),
		G: image.NewGray(image.Rect(0, 0, width, height)),
		B: image.NewGray(image.Rect(0, 0, width, height)),
		A: image.NewGray(image.Rect(0, 0, width, height)),
	}
}

// DrawShapes draws the given shapes onto the canvas.
func (c *Canvas) DrawShapes(shapes []Shape) error {
	for _, s := range shapes {
		switch shape := s.(type) {
		case Rectangle:
			if err := c.drawRectangle(shape); err != nil {
				return err
			}
		case Line:
			if err := c.drawLine(shape); err != nil {
				return err
			}
			// Add other shapes here if needed
		}
	}
	return nil
}

// Draw rectangle on canvas
func (c *Canvas) drawRectangle(r Rectangle) error {
	col, err := hexToRGBA(r.Color)
	if err != nil {
		return err
	}
	rect := image.Rect(r.X, r.Y, r.X+r.W, r.Y+r.H)
	draw.Draw(c.R, rect, &image.Uniform{C: color.Gray{Y: col.R}}, image.Point{}, draw.Src)
	draw.Draw(c.G, rect, &image.Uniform{C: color.Gray{Y: col.G}}, image.Point{}, draw.Src)
	draw.Draw(c.B, rect, &image.Uniform{C: color.Gray{Y: col.B}}, image.Point{}, draw.Src)
	draw.Draw(c.A, rect, &image.Uniform{C: color.Gray{Y: col.A}}, image.Point{}, draw.Src)
	return nil
}

// Draw line on canvas
func (c *Canvas) drawLine(l Line) error {
	col, err := hexToRGBA(l.Color)
	if err != nil {
		return err
	}

	var rect image.Rectangle
	half := l.Thickness / 2

	if l.X1 == l.X2 {
		y1, y2 := l.Y1, l.Y2
		if y1 > y2 {
			y1, y2 = y2, y1
		}

		rect = image.Rect(
			l.X1-half,
			y1,
			l.X1+half,
			y2,
		)

	} else if l.Y1 == l.Y2 {
		x1, x2 := l.X1, l.X2
		if x1 > x2 {
			x1, x2 = x2, x1
		}

		rect = image.Rect(
			x1,
			l.Y1-half,
			x2,
			l.Y1+half,
		)

	} else {
		return nil
	}

	draw.Draw(c.R, rect, &image.Uniform{C: color.Gray{Y: col.R}}, image.Point{}, draw.Src)
	draw.Draw(c.G, rect, &image.Uniform{C: color.Gray{Y: col.G}}, image.Point{}, draw.Src)
	draw.Draw(c.B, rect, &image.Uniform{C: color.Gray{Y: col.B}}, image.Point{}, draw.Src)
	draw.Draw(c.A, rect, &image.Uniform{C: color.Gray{Y: col.A}}, image.Point{}, draw.Src)
	return nil
}

// CalculateHashes calculates the SHA256 hash of each color channel.
func (c *Canvas) CalculateHashes() (map[string]string, error) {
	hashes := make(map[string]string)

	rHash, err := calculateHash(c.R.Pix)
	if err != nil {
		return nil, fmt.Errorf("failed to hash R channel: %w", err)
	}
	hashes["red"] = rHash

	gHash, err := calculateHash(c.G.Pix)
	if err != nil {
		return nil, fmt.Errorf("failed to hash G channel: %w", err)
	}
	hashes["green"] = gHash

	bHash, err := calculateHash(c.B.Pix)
	if err != nil {
		return nil, fmt.Errorf("failed to hash B channel: %w", err)
	}
	hashes["blue"] = bHash

	aHash, err := calculateHash(c.A.Pix)
	if err != nil {
		return nil, fmt.Errorf("failed to hash A channel: %w", err)
	}
	hashes["alpha"] = aHash

	return hashes, nil
}

// CalculateCombinedHash calculates the combined hash of the canvas.
func (c *Canvas) CalculateCombinedHash(results map[string]string) (string, error) {
	order := []string{"red", "green", "blue", "alpha"}

	combined := ""
	for _, k := range order {
		combined += results[k]
	}

	return calculateHash([]byte(combined))
}

func calculateHash(data []byte) (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write(data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func hexToRGBA(hexColor string) (color.RGBA, error) {
	var r, g, b byte
	_, err := fmt.Sscanf(hexColor, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("failed to parse hex color: %w", err)
	}
	return color.RGBA{R: r, G: g, B: b, A: 255}, nil
}

// PrintMatrices prints the pixel values of each color channel to the console.
func (c *Canvas) PrintMatrices() {
	fmt.Println("Red Channel:")
	printChannel(c.R)
	fmt.Println("Green Channel:")
	printChannel(c.G)
	fmt.Println("Blue Channel:")
	printChannel(c.B)
	fmt.Println("Alpha Channel:")
	printChannel(c.A)
}

func printChannel(channel *image.Gray) {
	bounds := channel.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			fmt.Printf("%3d ", channel.GrayAt(x, y).Y)
		}
		fmt.Println()
	}
}
