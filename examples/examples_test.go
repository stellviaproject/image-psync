package examples

import (
	"image"
	"image/color"
	"math/rand"
	"testing"
)

func TestGray(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, color.RGBA{
				R: uint8(rand.Int()),
				G: uint8(rand.Int()),
				B: uint8(rand.Int()),
				A: uint8(rand.Int()),
			})
		}
	}
	grayParallel := ConvertToGrayParallel(img, 10)
	gray := ConvertToGray(img)
	if !CompareImage(grayParallel, gray) {
		t.FailNow()
	}
}

func TestCountWindow(t *testing.T) {
	slice2d := make([][]bool, 100)
	for i := 0; i < 100; i++ {
		slice2d[i] = make([]bool, 100)
		for j := 0; j < 100; j++ {
			slice2d[i][j] = rand.Float64() > 0.5
		}
	}
	countParallel := CountWindowParallel(slice2d, 10, 5)
	count := CountWindow(slice2d, 10)
	if countParallel != count {
		t.FailNow()
	}
}

func TestSumSpectrum(t *testing.T) {
	slice2d := make([][]float64, 500)
	for i := 0; i < 500; i++ {
		slice2d[i] = make([]float64, 500)
		for j := 0; j < 500; j++ {
			slice2d[i][j] = rand.Float64()
		}
	}
	specParallel := MakeSumSpectrumParallel(slice2d, 10, 20)
	spec := MakeSumSpectrum(slice2d, 10)
	if !CompareFloat64(specParallel, spec) {
		t.FailNow()
	}
}

func TestSumAround(t *testing.T) {
	slice2d := make([][]float64, 500)
	for i := 0; i < 500; i++ {
		slice2d[i] = make([]float64, 500)
		for j := 0; j < 500; j++ {
			slice2d[i][j] = rand.Float64()
		}
	}
	aroundParallel := MakeSumAroundParallel(slice2d, 10, 60)
	around := MakeSumAround(slice2d, 10)
	if !CompareFloat64(aroundParallel, around) {
		t.FailNow()
	}
}

func CompareFloat64(src, dst [][]float64) bool {
	if len(src) != len(dst) {
		return false
	}
	for i := 0; i < len(src); i++ {
		if len(src[i]) != len(dst[i]) {
			return false
		}
		for j := 0; j < len(src[i]); j++ {
			if src[i][j] != dst[i][j] {
				return false
			}
		}
	}
	return true
}

func CompareImage(src, dst *image.Gray) bool {
	bdSrc, bdDst := src.Bounds(), dst.Bounds()
	if bdSrc.Dx() != bdDst.Dx() {
		return false
	}
	if bdSrc.Dy() != bdDst.Dy() {
		return false
	}
	width, height := bdSrc.Dx(), bdSrc.Dy()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			gSrc := src.GrayAt(x, y).Y
			gDst := dst.GrayAt(x, y).Y
			if gSrc != gDst {
				return false
			}
		}
	}
	return true
}
