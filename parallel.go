package psync

import "sync"

// It is use for passing to methods a function called on every region of image or slice in parallel.
// The values minX, minY are the inferior limit of region, and maxX, maxY are the superior limit.
type Action func(minX, minY, maxX, maxY int)

// This function receive the width, height of image or slice and divide (width / parallel is the size of every region, only last region can have a diferent size) it by regions of the same sizes using parallel, then run action on the region.
func ParallelRegionHorizontal(width, height, parallel int, action Action) {
	wdx := width / parallel
	wRegions := make([]int, parallel+1)
	for i := 0; i < parallel; i++ {
		wRegions[i] = i * wdx
	}
	wRegions[parallel] = width
	wg := sync.WaitGroup{}
	for i := 1; i < len(wRegions); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			action(wRegions[i-1], 0, wRegions[i], height)
		}(i)
	}
	wg.Wait()
}

// This function receive the width, height of image or slice and divide (height / parallel is the size of every region, only last region can have a diferent size) it by regions of the same sizes using parallel, then run action on the region.
func ParallelRegionVertical(width, height, parallel int, action Action) {
	hdy := height / parallel
	hRegions := make([]int, parallel+1)
	for i := 0; i < parallel; i++ {
		hRegions[i] = i * hdy
	}
	hRegions[parallel] = height
	wg := sync.WaitGroup{}
	for i := 1; i < len(hRegions); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			action(0, hRegions[i-1], width, hRegions[i])
		}(i)
	}
	wg.Wait()
}

// This function run action in parallel passing to action the coordinates of every window. The value of parallel control the number of gorutines.Note it's posible it doesn't cover every part of image, the number of windows in horizontal orientation will be number = width / window, but its division is an integer, then the last coordinates will be number * window and if width % window > 0 then an area of size = width % window won't be covered.
func ParallelWindow(width, height, window, parallel int, action Action) {
	prll := make(chan int, parallel)
	wg := sync.WaitGroup{}
	width, height = width-window, height-window
	for x := 0; x < width; x += window {
		for y := 0; y < height; y += window {
			wg.Add(1)
			go func(x, y int) {
				defer wg.Done()
				prll <- 0
				action(x, y, x+window, y+window)
				<-prll
			}(x, y)
		}
	}
	wg.Wait()
}

// This function run action in parallel for every pixel of image, the value of parallel control the number of gorutines. It is util when you want to process pixels around of every pixel.
func ParallelForEach(width, height, parallel int, action Action) {
	prll := make(chan int, parallel)
	wg := sync.WaitGroup{}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			wg.Add(1)
			go func(x, y int) {
				defer wg.Done()
				prll <- 0
				action(x, y, x+1, y+1)
				<-prll
			}(x, y)
		}
	}
	wg.Wait()
}

// This function create a grid where every cell have size of size, only the last cells in the corner may have a diferent size. It's create gorutines for every cell, and the value of parallel control the number of gorutines.
func ParallelGrid(width, height, size, parallel int, action Action) {

	numX := width / size
	gridX := make([]int, 0, numX+1)
	for i := 0; i < numX; i++ {
		gridX = append(gridX, i*size)
	}
	gridX = append(gridX, width)

	numY := height / size
	gridY := make([]int, 0, numY+1)
	for i := 0; i < numY; i++ {
		gridY = append(gridY, i*size)
	}
	gridY = append(gridY, height)

	prll := make(chan int, parallel)
	wg := sync.WaitGroup{}
	for i := 1; i < len(gridX); i++ {
		for j := 1; j < len(gridY); j++ {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				prll <- 0
				action(gridX[i-1], gridY[j-1], gridX[i], gridY[j])
				<-prll
			}(i, j)
		}
	}
	wg.Wait()
}
