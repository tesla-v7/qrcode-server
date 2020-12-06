package qrCode

import (
	"encoding/hex"
	"fmt"
	"github.com/skip2/go-qrcode"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"math"
	"os"
	"time"
)

type QrConf struct {
	ColorCenter string //"A0FF11"
	ColorEdge   string //"D1170F"
	//side        int
	LogoPath  string //"logo.png"
	PixRadius int
	Template  string
}

type rgba struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

type rgb struct {
	r uint8
	g uint8
	b uint8
}

type QrRender struct {
	colorGradientCenter *rgb
	colorGradientEdge   *rgb
	colorGradientDiff   *[3]float64
	gradientCenter      int
	gradientLength      float64

	side int

	pixRadius   int
	pixDiameter int
	pixel       *image.RGBA

	//logo          *image.Image
	logo          *image.RGBA
	backgroundPin *image.RGBA
	searchPattern *image.RGBA

	logoXY0, logoXY1 int
}

func (q *QrRender) Print() {
	fmt.Println(q.side)
}

func (q *QrRender) paintCircleGradient(x, y int) *color.RGBA {
	dx := x - q.gradientCenter
	dy := y - q.gradientCenter
	d := math.Sqrt(float64(dx*dx+dy*dy)) / q.gradientLength
	if d < 0 {
		d = 0
	}
	if d > 1 {
		d = 1
	}
	r := q.colorGradientCenter.r + byte(d*q.colorGradientDiff[0])
	g := q.colorGradientCenter.g + byte(d*q.colorGradientDiff[1])
	b := q.colorGradientCenter.b + byte(d*q.colorGradientDiff[2])

	return &color.RGBA{r, g, b, 255}
}

func (q *QrRender) paintBackgroundPix(sideLength int) {
	for x := 0; x < sideLength; x++ {
		for y := 0; y < sideLength; y++ {
			q.backgroundPin.Set(x, y, q.paintCircleGradient(x, y))
		}
	}
}

func TimeTrack(start time.Time) {
	timeRun := time.Since(start)
	fmt.Print("run time: ", timeRun, "\n")
}

func (q *QrRender) GetQrCodeImage(text *string) (*image.RGBA, error) {
	//t := time.Now()
	qr, err := qrcode.New(*text, qrcode.High)
	if err != nil {
		return nil, err
	}

	qrArray := qr.Bitmap()
	qrLength := len(qrArray[0]) - 4
	//fmt.Println("create QrBase64")
	//TimeTrack(t)
	//
	//t = time.Now()
	mask := image.NewRGBA(image.Rect(0, 0, q.side, q.side))
	qrCodeResult := image.NewRGBA(image.Rect(0, 0, q.side, q.side))
	radiusAddTwoPix := q.pixRadius + 2

	for x := 3; x < qrLength; x++ {
		for y := 3; y < qrLength; y++ {
			if (x < 11 && y < 11) || (x < 11 && y > qrLength-8) || (x > qrLength-8 && y < 11) {
				continue
			}
			if qrArray[y][x] {
				xStart := x - 3
				yStart := y - 3
				draw.DrawMask(
					mask,
					image.Rect(
						xStart*q.pixDiameter-q.pixRadius-xStart,
						yStart*q.pixDiameter-q.pixRadius-yStart,
						xStart*q.pixDiameter+radiusAddTwoPix-xStart,
						yStart*q.pixDiameter+radiusAddTwoPix-yStart),
					q.backgroundPin,
					image.ZP,
					q.pixel,
					image.ZP,
					draw.Over)
			}
		}
	}

	//fmt.Println("render mask QrBase64")
	//TimeTrack(t)
	//
	//t = time.Now()

	draw.DrawMask(qrCodeResult, qrCodeResult.Bounds(), q.backgroundPin, image.ZP, mask, image.ZP, draw.Over)

	xy0 := q.side - q.pixDiameter*8 + 8 // + q.pixRadius -7
	xy1 := q.pixDiameter*8 + 8

	draw.Draw(qrCodeResult, qrCodeResult.Bounds(), q.searchPattern, image.ZP, draw.Src)
	draw.Draw(qrCodeResult, image.Rect(0, xy0, xy1, q.side), q.searchPattern, image.ZP, draw.Src)
	draw.Draw(qrCodeResult, image.Rect(xy0, 0, q.side, xy1), q.searchPattern, image.ZP, draw.Src)

	draw.Draw(qrCodeResult, qrCodeResult.Bounds(), q.logo, image.ZP, draw.Over)
	//logoRec := image.Rect(q.logoXY0, q.logoXY0, q.logoXY1, q.logoXY1)
	//draw.BiLinear.Scale(qrCodeResult, logoRec.Bounds(), q.logo, q.logo.Bounds(), draw.Over, nil)

	//fmt.Println("render all QrBase64")
	//TimeTrack(t)
	return qrCodeResult, nil
}

func renderSearchPattern(searchPattern *image.RGBA, radius int, patternColor *rgb) {

	diameter := (radius << 1) - 2

	searchPatternTransparentStart := diameter + radius
	searchPatternTransparentEnd := diameter*6 + radius

	searchPatternTransparent := image.NewRGBA(image.Rect(searchPatternTransparentStart, searchPatternTransparentStart, searchPatternTransparentEnd, searchPatternTransparentEnd))

	searchPatternMinStart := diameter<<1 + radius
	searchPatterMinEnd := diameter*5 + radius
	searchPatternMin := image.NewRGBA(image.Rect(searchPatternMinStart, searchPatternMinStart, searchPatterMinEnd, searchPatterMinEnd))

	searchPatternColor := color.RGBA{patternColor.r, patternColor.g, patternColor.b, 255}
	searchPatternColorTransparent := color.RGBA{0, 0, 0, 0}

	draw.Draw(searchPattern, searchPattern.Bounds(), &image.Uniform{searchPatternColor}, image.ZP, draw.Src)
	draw.Draw(searchPattern, searchPatternTransparent.Bounds(), &image.Uniform{searchPatternColorTransparent}, image.ZP, draw.Src)
	draw.Draw(searchPattern, searchPatternMin.Bounds(), &image.Uniform{searchPatternColor}, image.ZP, draw.Over)
}

func strToRGB(strColor string) (*rgb, error) {
	rgbColor, err := hex.DecodeString(strColor)
	if err != nil {
		return nil, err
	}
	return &rgb{
		r: rgbColor[0],
		g: rgbColor[1],
		b: rgbColor[2],
	}, nil
}

func calcColorDiff(colorCenter *rgb, colorEdge *rgb) *[3]float64 {
	var result [3]float64
	result[0] = float64(colorEdge.r) - float64(colorCenter.r)
	result[1] = float64(colorEdge.g) - float64(colorCenter.g)
	result[2] = float64(colorEdge.b) - float64(colorCenter.b)
	//fmt.Println(colorCenter.b, colorEdge.b, result[2])
	return &result
}

func calcLogoPosition(diameter int, qrSide int) (int, int) {
	var logoRadius int = (diameter*8 - 1) >> 1
	var qrCenter int = qrSide >> 1
	var logoXY0, logoXY1 int = qrCenter - logoRadius, qrCenter + logoRadius
	return logoXY0, logoXY1
}

func New(qrConf *QrConf) (*QrRender, error) {

	var qrRender = QrRender{}
	var er error

	logoHwnd, err := os.Open(qrConf.LogoPath)

	defer logoHwnd.Close()

	if err != nil {
		return nil, err
	}

	logo, _, err := image.Decode(logoHwnd)
	if err != nil {
		return nil, err
	}

	//qrRender.logo = logo

	qrRender.colorGradientCenter, er = strToRGB(qrConf.ColorCenter)
	if er != nil {
		return nil, er
	}

	qrRender.colorGradientEdge, er = strToRGB(qrConf.ColorEdge)

	if er != nil {
		return nil, er
	}

	qr, err := qrcode.New(qrConf.Template, qrcode.High)
	if err != nil {
		return nil, err
	}

	qrArray := qr.Bitmap()
	qrLength := len(qrArray[0]) - 7
	side := qrLength*((qrConf.PixRadius<<1)-1) - qrLength

	qrRender.colorGradientDiff = calcColorDiff(qrRender.colorGradientCenter, qrRender.colorGradientEdge)
	qrRender.gradientCenter = side >> 1
	qrRender.pixRadius = qrConf.PixRadius
	qrRender.pixDiameter = (qrConf.PixRadius << 1) - 1
	qrRender.gradientLength = math.Sqrt(float64(side*side)*2)/2 - 40
	//fmt.Println(qrRender.gradientLength)
	qrRender.side = side
	//qrRender.logoXY0, qrRender.logoXY1 = calcLogoPosition(qrRender.pixDiameter, qrRender.side)
	logoXY0, logoXY1 := calcLogoPosition(qrRender.pixDiameter, qrRender.side)

	pixelDiameter := qrConf.PixRadius << 1
	searchPatternSize := (pixelDiameter-1)*7 + qrConf.PixRadius - 7

	qrRender.backgroundPin = image.NewRGBA(image.Rect(0, 0, side, side))
	qrRender.searchPattern = image.NewRGBA(image.Rect(qrConf.PixRadius-1, qrConf.PixRadius-1, searchPatternSize, searchPatternSize))
	qrRender.pixel = image.NewRGBA(image.Rect(0, 0, pixelDiameter+0, pixelDiameter+0))

	renderSearchPattern(qrRender.searchPattern, qrConf.PixRadius, qrRender.colorGradientEdge)
	draw.Draw(qrRender.pixel, qrRender.pixel.Bounds(), &Circle{image.Point{qrRender.pixRadius + 0, qrRender.pixRadius + 0}, qrRender.pixRadius + 0}, image.ZP, draw.Src)
	qrRender.paintBackgroundPix(side)

	qrRender.logo = image.NewRGBA(image.Rect(0, 0, side, side))
	logoRec := image.Rect(logoXY0, logoXY0, logoXY1, logoXY1)
	draw.BiLinear.Scale(qrRender.logo, logoRec.Bounds(), logo, logo.Bounds(), draw.Over, nil)

	return &qrRender, nil
}
