package _02testSuanfa

import (
	"fmt"
	"math"
	"testing"
)

func Test009Name(t *testing.T) {
	fmt.Println(isPowerOfThree(27))
}

func isPowerOfThree(n int) bool {
	maxNnm := math.MaxInt32
	temp := 3
	if n <= 0 || n > maxNnm {
		return false
	}
	for temp < maxNnm {
		if n == temp {
			return true
		}
		temp *= 3
	}
	return false
}
