package s3cr3ts4nt4

import (
	"fmt"
	"math/rand"
)

func DerangedIndices(n int) ([]int, error) {
	if n < 2 {
		return nil, fmt.Errorf("cannot create deranged indices for %d elements", n)
	}
	indices := make([]int, n)
	for i := 0; i < n; i++ {
		indices[i] = i
	}
	// Shift left by one to make sure nobody maps to themselves.
	indices = append(indices[1:], indices[0])
	rand.Shuffle(n, func(i, j int) {
		// Don't swap elements if it would result in an element ending up at
		// its original position.
		if indices[j] == i || indices[i] == j {
			return
		}
		indices[i], indices[j] = indices[j], indices[i]
	})

	return indices, nil
}
