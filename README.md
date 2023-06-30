# Documentation
## Introduction
This library was created for use easy parallelism on images, arrays and slices. It allows dividing an image or slice in regions with the same size and excecute an action over the region with a gorutine, or specifying the number of gorutines for dividing image horizontaly or verticaly, you can use a number of x and y to create a grid over image and running action in every cell in parallel too, it's posible to run an action on every pixel.

## Description of every object
This library is very simple at this moment and only has six objects, one is a <span style="color:blue">func</span> type and the others are functions.

The objects are:

| Names | Definition | Description of use |
|-------|------------|--------------------|
|Action | type Action func(minX, minY, maxX, maxY int) | It's a function called on every region of image or slice in parallel. The values **minX**, **minY** are the inferior limit of region, and **maxX**, **maxY** are the superior limit.|
|ParallelRegionHorizontal|func ParallelRegionHorizontal(width, height, parallel int, action Action)| This function receive the width, height of image or slice and divide (**width / parallel** is the size of every region, only last region can have a diferent size) it by regions of the same sizes using parallel, then run action on the region.|
|ParallelRegionVertical|func ParallelRegionVertical(width, height, parallel int, action Action)|This function receive the width, height of image or slice and divide (**height / parallel** is the size of every region, only last region can have a diferent size) it by regions of the same sizes using parallel, then run action on the region.|
|ParallelWindow|func ParallelWindow(width, height, window, parallel int, action Action)|This function run action in parallel passing to action the coordinates of every window. The value of **parallel** control the number of gorutines.Note it's posible it doesn't cover every part of image, the number of windows in horizontal orientation will be **number = width / window**, but its division is an integer, then the last coordinates will be **number * window** and if **width % window > 0** then an area of **size = width % window** won't be covered.|
|ParallelForEach|func ParallelForEach(width, height, parallel int, action Action)|This function run action in parallel for every pixel of image, the value of **parallel** control the number of gorutines. It is util when you want to process pixels around of every pixel.|
|ParallelGrid|func ParallelGrid(width, height, size, parallel int, action Action)|This function create a grid where every cell have size of **size**, only the last cells in the corner may have a diferent size. It's create gorutines for every cell, and the value of **parallel** control the number of gorutines.|

#### Nota:
Keep in mind that value of parallel must not exceed the size of image. At this moment this library doesn't test function parameters.

## Examples

Use parallelism for convert an image to gray scale:

```golang

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
```

Use parallelism for count pixels or matches of object (like in box-counting dimension algorithm of fractal geometry):

```golang
import psync "github.com/stellviaproject/image-psync"

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
```

Use parallelism for makeing operations in a regions of image (like evaluateing multifractal spectrum), in this example is the sum of every pixel in the cell:

```golang
import psync "github.com/stellviaproject/image-psync"

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
```

Use parallelism for makeing operations around a pixel:

```golang
import psync "github.com/stellviaproject/image-psync"

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
```

### License

This code is available under the MIT License, which is a free and open source license. The terms and conditions of the license can be found at [enlace](https://github.com/stellviaproject/image-psync/LICENSE).

By using this code, you are required to follow the terms and conditions of the MIT License, which include including a copy of the license in the software and attributing the corresponding credits. If you have any questions about the license or how to use the code, please do not hesitate to contact us.
