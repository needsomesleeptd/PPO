package bboxes_utils

import (
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type BoundingBox struct {
	XMin, YMin, XMax, YMax int
}

func DrawText(img *image.RGBA, x, y int, text string) {

	point := fixed.Point26_6{fixed.I(x), fixed.I(y)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{200, 100, 0, 255}), // Configure the text color here
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)
}

// DrawBoundingBoxes draws bounding boxes on the image
func DrawBoundingBoxes(img *image.RGBA, boxes []BoundingBox, color color.RGBA) {
	for _, box := range boxes {
		drawRectangle(img, box, color)

	}
}

func drawRectangle(img *image.RGBA, box BoundingBox, color color.RGBA) {
	// Draw top and bottom lines
	drawLine(img, box.XMin, box.YMin, box.XMax, box.YMin, color)
	drawLine(img, box.XMin, box.YMax, box.XMax, box.YMax, color)

	// Draw left and right lines
	drawLine(img, box.XMin, box.YMin, box.XMin, box.YMax, color)
	drawLine(img, box.XMax, box.YMin, box.XMax, box.YMax, color)
}

func drawLine(img *image.RGBA, x0, y0, x1, y1 int, color color.RGBA) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx - dy

	for {
		img.Set(x0, y0, color)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
