package psync

import "sync"

type Action func(minX, minY, maxX, maxY int)

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
