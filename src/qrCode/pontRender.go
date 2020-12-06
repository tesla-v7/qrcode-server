package qrCode

import (
	"image"
	"image/color"
	"math"
)

type Circle struct {
	p image.Point
	r int
}

func (c *Circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *Circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *Circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X), float64(y-c.p.Y), float64(c.r)
	//xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)

	kP := rr - (math.Sqrt(xx*xx + yy*yy))
	if kP <= 0 {
		return color.RGBA{0, 0, 0, 0}
	}

	ka := 255 - (255 * kP)
	if ka < 0 {
		return color.RGBA{255, 255, 255, 255}
	}

	colorGrad := uint8(255 - ka)
	return color.RGBA{colorGrad, colorGrad, colorGrad, colorGrad}
}
