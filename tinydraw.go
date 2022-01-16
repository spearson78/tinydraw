package tinydraw // import "tinygo.org/x/tinydraw"

import (
	"image/color"

	"errors"

	"tinygo.org/x/drivers"
)

type DisplayerEx interface {
	drivers.Displayer
	FillRect(x, y, w, h int16, c color.RGBA)
}

type fillRectAdapter struct {
	drivers.Displayer
}

func (adapter *fillRectAdapter) FillRect(x, y, w, h int16, c color.RGBA) {
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			adapter.Displayer.SetPixel(px, py, c)
		}
	}
}

// Line draws a line between two points
func Line(display drivers.Displayer, x0 int16, y0 int16, x1 int16, y1 int16, color color.RGBA) {
	adapter := fillRectAdapter{display}
	LineEx(&adapter, x0, y0, x1, y1, color)
}

func LineEx(display DisplayerEx, x0 int16, y0 int16, x1 int16, y1 int16, color color.RGBA) {
	if x0 == x1 {
		if y0 > y1 {
			y0, y1 = y1, y0
		}
		display.FillRect(x0, y0, 1, (y1-y0)+1, color)
	} else if y0 == y1 {
		if x0 > x1 {
			x0, x1 = x1, x0
		}

		display.FillRect(x0, y0, (x1-x0)+1, 1, color)
	} else { // Bresenham
		dx := x1 - x0
		if dx < 0 {
			dx = -dx
		}
		dy := y1 - y0
		if dy < 0 {
			dy = -dy
		}
		steep := dy > dx
		if steep {
			x0, x1, y0, y1 = y0, y1, x0, x1
		}
		if x0 > x1 {
			x0, x1, y0, y1 = x1, x0, y1, y0
		}
		dx = x1 - x0
		dy = y1 - y0
		ystep := int16(1)
		if dy < 0 {
			dy = -dy
			ystep = -1
		}
		err := dx / 2

		startx := x0

		for ; x0 <= x1; x0++ {
			err -= dy
			if err < 0 {

				if steep {
					display.FillRect(y0, startx, 1, (x0-startx)+1, color)
				} else {
					display.FillRect(startx, y0, (x0-startx)+1, 1, color)
				}

				startx = x0 + 1

				y0 += ystep
				err += dx
			}
		}

		if x0 != startx {
			if steep {
				display.FillRect(y0, startx, 1, (x0 - startx), color)
			} else {
				display.FillRect(startx, y0, (x0 - startx), 1, color)
			}
		}
	}
}

// Rectangle draws a rectangle given a point, width and height
func Rectangle(display drivers.Displayer, x int16, y int16, w int16, h int16, color color.RGBA) error {
	adapter := fillRectAdapter{display}
	return RectangleEx(&adapter, x, y, w, h, color)
}

func RectangleEx(display DisplayerEx, x int16, y int16, w int16, h int16, color color.RGBA) error {
	if w <= 0 || h <= 0 {
		return errors.New("empty rectangle")
	}
	LineEx(display, x, y, x+w-1, y, color)
	LineEx(display, x, y, x, y+h-1, color)
	LineEx(display, x+w-1, y, x+w-1, y+h-1, color)
	LineEx(display, x, y+h-1, x+w-1, y+h-1, color)
	return nil
}

// FilledRectangle draws a filled rectangle given a point, width and height
func FilledRectangle(display drivers.Displayer, x int16, y int16, w int16, h int16, color color.RGBA) error {
	adapter := fillRectAdapter{display}
	return FilledRectangleEx(&adapter, x, y, w, h, color)
}

func FilledRectangleEx(display DisplayerEx, x int16, y int16, w int16, h int16, color color.RGBA) error {
	if w <= 0 || h <= 0 {
		return errors.New("empty rectangle")
	}
	display.FillRect(x, y, w, h, color)
	return nil
}

// Circle draws a circle given a point and radius
func Circle(display drivers.Displayer, x0 int16, y0 int16, r int16, color color.RGBA) {
	f := 1 - r
	ddfx := int16(1)
	ddfy := -2 * r
	x := int16(0)
	y := r
	display.SetPixel(x0, y0+r, color)
	display.SetPixel(x0, y0-r, color)
	display.SetPixel(x0+r, y0, color)
	display.SetPixel(x0-r, y0, color)
	for x < y {
		if f >= 0 {
			y--
			ddfy += 2
			f += ddfy
		}
		x++
		ddfx += 2
		f += ddfx

		display.SetPixel(x0+x, y0+y, color)
		display.SetPixel(x0-x, y0+y, color)
		display.SetPixel(x0+x, y0-y, color)
		display.SetPixel(x0-x, y0-y, color)
		display.SetPixel(x0+y, y0+x, color)
		display.SetPixel(x0-y, y0+x, color)
		display.SetPixel(x0+y, y0-x, color)
		display.SetPixel(x0-y, y0-x, color)
	}
}

// FilledCircle draws a filled circle given a point and radius
func FilledCircle(display drivers.Displayer, x0 int16, y0 int16, r int16, color color.RGBA) {
	adapter := fillRectAdapter{display}
	FilledCircleEx(&adapter, x0, y0, r, color)
}

func FilledCircleEx(display DisplayerEx, x0 int16, y0 int16, r int16, color color.RGBA) {
	f := 1 - r
	ddfx := int16(1)
	ddfy := -2 * r
	x := int16(0)
	y := r
	LineEx(display, x0, y0-r, x0, y0+r, color)
	for x < y {
		if f >= 0 {
			y--
			ddfy += 2
			f += ddfy
		}
		x++
		ddfx += 2
		f += ddfx

		LineEx(display, x0+x, y0-y, x0+x, y0+y, color)
		LineEx(display, x0+y, y0-x, x0+y, y0+x, color)
		LineEx(display, x0-x, y0-y, x0-x, y0+y, color)
		LineEx(display, x0-y, y0-x, x0-y, y0+x, color)
	}
}

// Triangle draws a triangle given three points
func Triangle(display drivers.Displayer, x0 int16, y0 int16, x1 int16, y1 int16, x2 int16, y2 int16, color color.RGBA) {
	adapter := fillRectAdapter{display}
	TriangleEx(&adapter, x0, y0, x1, y1, x2, y2, color)
}

func TriangleEx(display DisplayerEx, x0 int16, y0 int16, x1 int16, y1 int16, x2 int16, y2 int16, color color.RGBA) {
	LineEx(display, x0, y0, x1, y1, color)
	LineEx(display, x0, y0, x2, y2, color)
	LineEx(display, x1, y1, x2, y2, color)
}

// FilledTriangle draws a filled triangle given three points
func FilledTriangle(display drivers.Displayer, x0 int16, y0 int16, x1 int16, y1 int16, x2 int16, y2 int16, color color.RGBA) {
	adapter := fillRectAdapter{display}
	FilledTriangleEx(&adapter, x0, y0, x1, y1, x2, y2, color)
}

func FilledTriangleEx(display DisplayerEx, x0 int16, y0 int16, x1 int16, y1 int16, x2 int16, y2 int16, color color.RGBA) {
	if y0 > y1 {
		x0, y0, x1, y1 = x1, y1, x0, y0
	}
	if y1 > y2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}
	if y0 > y1 {
		x0, y0, x1, y1 = x1, y1, x0, y0
	}

	if y0 == y2 { // y0 = y1 = y2 : it's a line
		a := x0
		b := x0
		if x1 < a {
			a = x1
		} else if x1 > b {
			b = x1
		}
		if x2 < a {
			a = x2
		} else if x2 > b {
			b = x2
		}
		LineEx(display, a, y0, b, y0, color)
		return
	}

	dx01 := x1 - x0
	dy01 := y1 - y0
	dx02 := x2 - x0
	dy02 := y2 - y0
	dx12 := x2 - x1
	dy12 := y2 - y1

	sa := int16(0)
	sb := int16(0)
	a := int16(0)
	b := int16(0)

	last := y1 - 1
	if y1 == y2 {
		last = y1
	}

	y := y0
	for ; y <= last; y++ {
		a = x0 + sa/dy01
		b = x0 + sb/dy02
		sa += dx01
		sb += dx02
		LineEx(display, a, y, b, y, color)
	}

	sa = dx12 * (y - y1)
	sb = dx02 * (y - y0)

	for ; y <= y2; y++ {
		a = x1 + sa/dy12
		b = x0 + sb/dy02
		sa += dx12
		sb += dx02
		LineEx(display, a, y, b, y, color)
	}
}
