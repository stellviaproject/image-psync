package psync

import "testing"

func TestParallelGrid(t *testing.T) {
	ParallelGrid(100, 100, 50, 10, func(minX, minY, maxX, maxY int) {})
	ParallelGrid(130, 130, 50, 10, func(minX, minY, maxX, maxY int) {})
}
