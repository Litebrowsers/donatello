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

package canvas_tasks

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Shape interface defines the common behavior for all shapes.
type Shape interface {
	Encode() string
}

// Rectangle represents a rectangle shape.
type Rectangle struct {
	Color string
	W, H  int
	X, Y  int
}

// Encode returns the encoded string for a Rectangle.
func (r Rectangle) Encode() string {
	return fmt.Sprintf("R:%s:%d:%d:%d:%d", r.Color, r.W, r.H, r.X, r.Y)
}

// BoundingBox returns the bounding box of the Rectangle.
func (r Rectangle) BoundingBox() Rect {
	return Rect{r.X, r.Y, r.X + r.W, r.Y + r.H}
}

// Circle represents a circle shape.
type Circle struct {
	Color string
	R     int
	X, Y  int
}

// Encode returns the encoded string for a Circle.
func (c Circle) Encode() string {
	return fmt.Sprintf("C:%s:%d:%d:%d", c.Color, c.R, c.X, c.Y)
}

// Triangle represents a triangle shape.
type Triangle struct {
	Color                  string
	X1, Y1, X2, Y2, X3, Y3 int
}

// Encode returns the encoded string for a Triangle.
func (t Triangle) Encode() string {
	return fmt.Sprintf("T:%s:%d:%d:%d:%d:%d:%d", t.Color, t.X1, t.Y1, t.X2, t.Y2, t.X3, t.Y3)
}

// Line represents a line shape.

type Line struct {
	Color          string
	X1, Y1, X2, Y2 int
	Thickness      int
}

// Encode returns the encoded string for a Line.
func (l Line) Encode() string {
	return fmt.Sprintf("L:%s:%d:%d:%d:%d:%d", l.Color, l.X1, l.Y1, l.X2, l.Y2, l.Thickness)
}

// BoundingBox returns the bounding box of the Line.
func (l Line) BoundingBox() Rect {
	minX := min(l.X1, l.X2)
	maxX := max(l.X1, l.X2)
	minY := min(l.Y1, l.Y2)
	maxY := max(l.Y1, l.Y2)
	return Rect{minX, minY, maxX + l.Thickness, maxY + l.Thickness} // Approximation for thick lines
}

// Ellipse represents an ellipse shape.
type Ellipse struct {
	Color  string
	RX, RY int
	X, Y   int
}

// Encode returns the encoded string for an Ellipse.
func (e Ellipse) Encode() string {
	return fmt.Sprintf("E:%s:%d:%d:%d:%d", e.Color, e.RX, e.RY, e.X, e.Y)
}

// Rect represents an axis-aligned bounding box.
type Rect struct {
	MinX, MinY, MaxX, MaxY int
}

// Helper functions for min/max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// TaskGenerator generates encoded task strings from a list of shapes.
type TaskGenerator struct {
	shapes []Shape
}

// NewTaskGenerator creates a new TaskGenerator with the given shapes.
func NewTaskGenerator(shapes ...Shape) *TaskGenerator {
	return &TaskGenerator{
		shapes: shapes,
	}
}

// AddShape adds a shape to the generator.
func (tg *TaskGenerator) AddShape(s Shape) {
	tg.shapes = append(tg.shapes, s)
}

// GenerateTask returns a combined encoded string of all shapes.
func (tg *TaskGenerator) GenerateTask() string {
	var encodedShapes []string
	for _, s := range tg.shapes {
		encodedShapes = append(encodedShapes, s.Encode())
	}
	return strings.Join(encodedShapes, ";")
}

// generateRandomColor generates a random 6-digit hexadecimal color string.
func generateRandomColor() string {
	return fmt.Sprintf("%06X", rand.Intn(0xFFFFFF+1))
}

// Overlaps checks if two Rectangles overlap.
func (r1 Rect) Overlaps(r2 Rect) bool {
	return r1.MinX < r2.MaxX && r1.MaxX > r2.MinX &&
		r1.MinY < r2.MaxY && r1.MaxY > r2.MinY
}

// checkOverlap checks if a new shape overlaps with any existing shapes.
func checkOverlap(newShape Shape, existingShapes []Shape) bool {
	newRect := getBoundingBox(newShape)
	for _, existingShape := range existingShapes {
		existingRect := getBoundingBox(existingShape)
		if newRect.Overlaps(existingRect) {
			return true
		}
	}
	return false
}

// getBoundingBox is a helper to get the bounding box for any Shape.
func getBoundingBox(s Shape) Rect {
	switch shape := s.(type) {
	case Rectangle:
		return shape.BoundingBox()
	case Line:
		return shape.BoundingBox()
	default:
		return Rect{}
	}
}

// GenerateRandomShapes generates a slice of random shapes.
func GenerateRandomShapes(canvasSize int, count int) []Shape {
	if count <= 0 {
		return nil
	}

	shapes := make([]Shape, count)
	for i := 0; i < count; i++ {
		switch rand.Intn(5) { // 0: Rectangle, 1: Circle, 2: Triangle, 3: Line, 4: Ellipse
		case 0:
			shapes[i] = Rectangle{
				Color: generateRandomColor(),
				W:     rand.Intn(canvasSize/2) + 1,
				H:     rand.Intn(canvasSize/2) + 1,
				X:     rand.Intn(canvasSize),
				Y:     rand.Intn(canvasSize),
			}
		case 1:
			shapes[i] = Circle{
				Color: generateRandomColor(),
				R:     rand.Intn(canvasSize/4) + 1,
				X:     rand.Intn(canvasSize),
				Y:     rand.Intn(canvasSize),
			}
		case 2:
			shapes[i] = Triangle{
				Color: generateRandomColor(),
				X1:    rand.Intn(canvasSize),
				Y1:    rand.Intn(canvasSize),
				X2:    rand.Intn(canvasSize),
				Y2:    rand.Intn(canvasSize),
				X3:    rand.Intn(canvasSize),
				Y3:    rand.Intn(canvasSize),
			}
		case 3:
			shapes[i] = Line{
				Color: generateRandomColor(),
				X1:    rand.Intn(canvasSize),
				Y1:    rand.Intn(canvasSize),
				X2:    rand.Intn(canvasSize),
				Y2:    rand.Intn(canvasSize),
			}
		case 4:
			shapes[i] = Ellipse{
				Color: generateRandomColor(),
				RX:    rand.Intn(canvasSize/2) + 1,
				RY:    rand.Intn(canvasSize/2) + 1,
				X:     rand.Intn(canvasSize),
				Y:     rand.Intn(canvasSize),
			}
		}
	}
	return shapes
}

// GenerateRandomEvenSizedSquares generates a slice of random square shapes with even side lengths.
func GenerateRandomEvenSizedPrimitives(canvasSize int, count int) []Shape {
	if count <= 0 {
		return nil
	}

	primitives := make([]Shape, 0, count) // Use 0 for initial length, count for capacity
	maxRetries := 100                     // Limit retries to prevent infinite loops

	for i := 0; i < count; i++ {
		retries := 0
		for retries < maxRetries {
			var newShape Shape
			switch rand.Intn(2) { // 0: Even-sided Square, 1: Line
			case 0:
				side := (rand.Intn(5) + 1) * 2
				newShape = Rectangle{
					Color: generateRandomColor(),
					W:     side,
					H:     side,
					X:     rand.Intn(canvasSize - side),
					Y:     rand.Intn(canvasSize - side),
				}
			case 1:
				thickness := (rand.Intn(2) + 1) * 2 // 2 or 4
				if rand.Intn(2) == 0 {
					x := rand.Intn(canvasSize - thickness)
					y1 := rand.Intn(canvasSize)
					y2 := rand.Intn(canvasSize)
					newShape = Line{
						Color:     generateRandomColor(),
						X1:        x,
						Y1:        y1,
						X2:        x,
						Y2:        y2,
						Thickness: thickness,
					}
				} else {
					y := rand.Intn(canvasSize - thickness)
					x1 := rand.Intn(canvasSize)
					x2 := rand.Intn(canvasSize)
					newShape = Line{
						Color:     generateRandomColor(),
						X1:        x1,
						Y1:        y,
						X2:        x2,
						Y2:        y,
						Thickness: thickness,
					}
				}
			}

			if !checkOverlap(newShape, primitives) {
				primitives = append(primitives, newShape)
				break
			}
			retries++
		}
		if retries == maxRetries {
			fmt.Printf("Warning: Could not place shape %d after %d retries. Canvas might be full or shapes too large.\n", i, maxRetries)
		}
	}
	return primitives
}
