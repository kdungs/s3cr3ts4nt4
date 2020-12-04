package s3cr3ts4nt4

import (
	"testing"

	"github.com/matryer/is"
)

func TestDerangedIndices(t *testing.T) {
	is := is.New(t)

	goodCases := []int{
		2,
		5,
		10,
		23,
		100,
		1000,
	}
	for _, n := range goodCases {
		idx, err := DerangedIndices(n)
		is.NoErr(err)

		// Make sure all numbers show up.
		found := make(map[int]bool, n)
		for _, i := range idx {
			found[i] = true
		}
		for i := 0; i < n; i++ {
			f, ok := found[i]
			is.True(ok)
			is.True(f)
		}

		// Make sure nobody gets mapped to themselves.
		for from, to := range idx {
			is.True(from != to)
		}
	}

	failureCases := []int{
		-1,
		0,
		1,
	}
	for _, n := range failureCases {
		_, err := DerangedIndices(n)
		is.True(err != nil)
	}
}
