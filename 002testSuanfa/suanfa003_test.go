package _02testSuanfa

import (
	"fmt"
	"testing"
)

func TestSuanfa003(t *testing.T) {
	fmt.Println(checkTwoChessboards("a1", "b2"))
}

func checkTwoChessboards(coordinate1 string, coordinate2 string) bool {
	//byt := coordinate1[0]
	a1 := (coordinate1[0] - 'a' + coordinate1[1] - '1') % 2
	a2 := (coordinate2[0] - 'a' + coordinate2[1] - '1') % 2
	return a1 == a2
}

// 两数之和
func temp112(nums []int, target int) []int {
	map1 := make(map[int]int, len(nums))
	for i := 0; i < len(nums); i++ {
		if j, ok := map1[target-nums[i]]; ok {
			return []int{j, i}
		}
		map1[nums[i]] = i
	}
	return nil
}

// 盛水最多的容器
func maxArea(height []int) int {
	var left, right, maxArea1 int = 0, len(height) - 1, 0
	for left < right {
		area := (right - left) * min(height[left], height[right])
		maxArea1 = max(maxArea1, area)
		if height[left] < height[right] {
			left++
		} else {
			right--
		}
	}
	return maxArea1
}

//接雨水
