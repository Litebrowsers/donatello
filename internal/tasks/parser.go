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
	"fmt"
	"strconv"
	"strings"
)

// ParseTask parses an encoded task string and returns a slice of shapes.
func ParseTask(task string) ([]Shape, error) {
	encodedShapes := strings.Split(task, ";")
	shapes := make([]Shape, 0, len(encodedShapes))

	for _, encodedShape := range encodedShapes {
		parts := strings.Split(encodedShape, ":")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid shape format: %s", encodedShape)
		}

		shapeType := parts[0]
		switch shapeType {
		case "R":
			shape, err := parseRectangle(parts)
			if err != nil {
				return nil, err
			}
			shapes = append(shapes, shape)
		case "C":
			shape, err := parseCircle(parts)
			if err != nil {
				return nil, err
			}
			shapes = append(shapes, shape)
		case "T":
			shape, err := parseTriangle(parts)
			if err != nil {
				return nil, err
			}
			shapes = append(shapes, shape)
		case "L":
			shape, err := parseLine(parts)
			if err != nil {
				return nil, err
			}
			shapes = append(shapes, shape)
		case "E":
			shape, err := parseEllipse(parts)
			if err != nil {
				return nil, err
			}
			shapes = append(shapes, shape)
		default:
			return nil, fmt.Errorf("unknown shape type: %s", shapeType)
		}
	}

	return shapes, nil
}

func parseRectangle(parts []string) (Rectangle, error) {
	if len(parts) != 6 {
		return Rectangle{}, fmt.Errorf("invalid rectangle format: %v", parts)
	}
	w, err := strconv.Atoi(parts[2])
	if err != nil {
		return Rectangle{}, err
	}
	h, err := strconv.Atoi(parts[3])
	if err != nil {
		return Rectangle{}, err
	}
	x, err := strconv.Atoi(parts[4])
	if err != nil {
		return Rectangle{}, err
	}
	y, err := strconv.Atoi(parts[5])
	if err != nil {
		return Rectangle{}, err
	}
	return Rectangle{Color: parts[1], W: w, H: h, X: x, Y: y}, nil
}

func parseCircle(parts []string) (Circle, error) {
	if len(parts) != 5 {
		return Circle{}, fmt.Errorf("invalid circle format: %v", parts)
	}
	r, err := strconv.Atoi(parts[2])
	if err != nil {
		return Circle{}, err
	}
	x, err := strconv.Atoi(parts[3])
	if err != nil {
		return Circle{}, err
	}
	y, err := strconv.Atoi(parts[4])
	if err != nil {
		return Circle{}, err
	}
	return Circle{Color: parts[1], R: r, X: x, Y: y}, nil
}

func parseTriangle(parts []string) (Triangle, error) {
	if len(parts) != 8 {
		return Triangle{}, fmt.Errorf("invalid triangle format: %v", parts)
	}
	x1, err := strconv.Atoi(parts[2])
	if err != nil {
		return Triangle{}, err
	}
	y1, err := strconv.Atoi(parts[3])
	if err != nil {
		return Triangle{}, err
	}
	x2, err := strconv.Atoi(parts[4])
	if err != nil {
		return Triangle{}, err
	}
	y2, err := strconv.Atoi(parts[5])
	if err != nil {
		return Triangle{}, err
	}
	x3, err := strconv.Atoi(parts[6])
	if err != nil {
		return Triangle{}, err
	}
	y3, err := strconv.Atoi(parts[7])
	if err != nil {
		return Triangle{}, err
	}
	return Triangle{Color: parts[1], X1: x1, Y1: y1, X2: x2, Y2: y2, X3: x3, Y3: y3}, nil
}

func parseLine(parts []string) (Line, error) {
	if len(parts) != 7 {
		return Line{}, fmt.Errorf("invalid line format: %v", parts)
	}
	x1, err := strconv.Atoi(parts[2])
	if err != nil {
		return Line{}, err
	}
	y1, err := strconv.Atoi(parts[3])
	if err != nil {
		return Line{}, err
	}
	x2, err := strconv.Atoi(parts[4])
	if err != nil {
		return Line{}, err
	}
	y2, err := strconv.Atoi(parts[5])
	if err != nil {
		return Line{}, err
	}
	thickness, err := strconv.Atoi(parts[6])
	if err != nil {
		return Line{}, err
	}
	return Line{Color: parts[1], X1: x1, Y1: y1, X2: x2, Y2: y2, Thickness: thickness}, nil
}

func parseEllipse(parts []string) (Ellipse, error) {
	if len(parts) != 6 {
		return Ellipse{}, fmt.Errorf("invalid ellipse format: %v", parts)
	}
	rx, err := strconv.Atoi(parts[2])
	if err != nil {
		return Ellipse{}, err
	}
	ry, err := strconv.Atoi(parts[3])
	if err != nil {
		return Ellipse{}, err
	}
	x, err := strconv.Atoi(parts[4])
	if err != nil {
		return Ellipse{}, err
	}
	y, err := strconv.Atoi(parts[5])
	if err != nil {
		return Ellipse{}, err
	}
	return Ellipse{Color: parts[1], RX: rx, RY: ry, X: x, Y: y}, nil
}
