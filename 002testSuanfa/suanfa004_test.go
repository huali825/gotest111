package _02testSuanfa

import (
	"testing"
)

func TestSuanfa004(t *testing.T) {
	slice1 := make([]int, 2)
	//slice2 := []int{1, 2, 3, 4, 5}
	for i := 0; i < 100; i++ {
		slice1 = append(slice1, i)
	}

	//nextslicecap(100, 2)

	//nums2 := []int{2, 7, 33, 44}
	//temp := twoNums(nums2, 9)
	//fmt.Println(temp)
}

func twoNums(nums []int, target int) []int {
	map1 := make(map[int]int, len(nums))
	for i := 0; i < len(nums); i++ {
		if j, ok := map1[target-nums[i]]; ok {
			return []int{i, j}
		} else {
			map1[nums[i]] = i
		}
	}
	return nil
}

func nextslicecap(newLen, oldCap int) int {
	newcap := oldCap
	doublecap := newcap + newcap
	if newLen > doublecap {
		return newLen
	}

	const threshold = 256
	if oldCap < threshold {
		return doublecap
	}
	for {
		// Transition from growing 2x for small slices
		// to growing 1.25x for large slices. This formula
		// gives a smooth-ish transition between the two.
		newcap += (newcap + 3*threshold) >> 2

		// We need to check `newcap >= newLen` and whether `newcap` overflowed.
		// newLen is guaranteed to be larger than zero, hence
		// when newcap overflows then `uint(newcap) > uint(newLen)`.
		// This allows to check for both with the same comparison.
		if uint(newcap) >= uint(newLen) {
			break
		}
	}

	// Set newcap to the requested cap when
	// the newcap calculation overflowed.
	if newcap <= 0 {
		return newLen
	}
	return newcap
}
