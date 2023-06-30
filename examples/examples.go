package examples

import (
	"image"
	"image/color"

	psync "github.com/stellviaproject/image-psync"
)

func ConvertToGrayParallel(img image.Image, parallel int) *image.Gray {
	bd := img.Bounds()
	gray := image.NewGray(bd)
	psync.ParallelRegionHorizontal(bd.Dx(), bd.Dy(), parallel, func(minX, minY, maxX, maxY int) {
		for x := minX; x < maxX; x++ {
			for y := minY; y < maxY; y++ {
				r, g, b, _ := img.At(x, y).RGBA()
				grayColor := uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
				gray.SetGray(x, y, color.Gray{Y: grayColor})
			}
		}
	})
	return gray
}

func ConvertToGray(img image.Image) *image.Gray {
	bd := img.Bounds()
	gray := image.NewGray(bd)
	width, height := bd.Dx(), bd.Dy()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			grayColor := uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
			gray.SetGray(x, y, color.Gray{Y: grayColor})
		}
	}
	return gray
}

func CountWindowParallel(slice2d [][]bool, window int, parallel int) int {
	width, height := len(slice2d[0]), len(slice2d)
	count := 0
	mutex := make(chan int, 1)
	psync.ParallelWindow(width, height, window, parallel, func(minX, minY, maxX, maxY int) {
		for x := minX; x < maxX; x++ {
			for y := minY; y < maxY; y++ {
				if slice2d[y][x] {
					mutex <- 0
					count++
					<-mutex
					return
				}
			}
		}
	})
	return count
}

func CountWindow(slice2d [][]bool, window int) int {
	width, height := len(slice2d[0]), len(slice2d)
	count := 0
	for minX := 0; minX < width; minX += window {
		for minY := 0; minY < height; minY += window {
			maxX, maxY := minX+window, minY+window
			if maxX < width && maxY < height {
			R:
				for x := minX; x < maxX; x++ {
					for y := minY; y < maxY; y++ {
						if slice2d[y][x] {
							count++
							break R
						}
					}
				}
			}
		}
	}
	return count
}

func MakeSumSpectrumParallel(slice2d [][]float64, size int, parallel int) [][]float64 {
	width, height := len(slice2d[0]), len(slice2d)
	sum := make([][]float64, height)
	for i := 0; i < len(slice2d); i++ {
		sum[i] = make([]float64, width)
	}
	psync.ParallelGrid(width, height, size, parallel, func(minX, minY, maxX, maxY int) {
		valueSum := 0.0
		for x := minX; x < maxX; x++ {
			for y := minY; y < maxY; y++ {
				valueSum += slice2d[y][x]
			}
		}
		for x := minX; x < maxX; x++ {
			for y := minY; y < maxY; y++ {
				sum[y][x] = valueSum
			}
		}
	})
	return sum
}

func MakeSumSpectrum(slice2d [][]float64, size int) [][]float64 {
	width, height := len(slice2d[0]), len(slice2d)
	sum := make([][]float64, height)
	for i := 0; i < len(slice2d); i++ {
		sum[i] = make([]float64, width)
	}
	for minX := 0; minX < width; minX += size {
		for minY := 0; minY < height; minY += size {
			maxX, maxY := minX+size, minY+size
			if maxX >= width {
				maxX = width
			}
			if maxY >= height {
				maxY = height
			}
			valueSum := 0.0
			for x := minX; x < maxX; x++ {
				for y := minY; y < maxY; y++ {
					valueSum += slice2d[y][x]
				}
			}
			for x := minX; x < maxX; x++ {
				for y := minY; y < maxY; y++ {
					sum[y][x] = valueSum
				}
			}
		}
	}
	return sum
}

func MakeSumAroundParallel(slice2d [][]float64, size, parallel int) [][]float64 {
	width, height := len(slice2d[0]), len(slice2d)
	around := make([][]float64, height)
	for i := 0; i < len(slice2d); i++ {
		around[i] = make([]float64, width)
	}
	psync.ParallelForEach(width, height, parallel, func(minX, minY, maxX, maxY int) {
		valueSum := 0.0
		for x := -size; x < size; x++ {
			for y := -size; y < size; y++ {
				x0, y0 := minX+x, minY+y
				if x0 < 0 || y0 < 0 || x0 >= width || y0 >= height {
					continue
				}
				valueSum += slice2d[y0][x0]
			}
		}
		around[minX][minY] = valueSum
	})
	return around
}

func MakeSumAround(slice2d [][]float64, ratio int) [][]float64 {
	width, height := len(slice2d[0]), len(slice2d)
	around := make([][]float64, height)
	for i := 0; i < len(slice2d); i++ {
		around[i] = make([]float64, width)
	}
	for minX := 0; minX < width; minX++ {
		for minY := 0; minY < height; minY++ {
			valueSum := 0.0
			for x := -ratio; x < ratio; x++ {
				for y := -ratio; y < ratio; y++ {
					x0, y0 := minX+x, minY+y
					if x0 < 0 || y0 < 0 || x0 >= width || y0 >= height {
						continue
					}
					valueSum += slice2d[y0][x0]
				}
			}
			around[minX][minY] = valueSum
		}
	}
	return around
}
