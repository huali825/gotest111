package testSuanfa

import "testing"

func TestLiangshuzhihe(t *testing.T) {
	nums := []int{2, 7, 11, 15}
	t.Log(twoSum(nums, 9))
}

func twoSum(nums []int, target int) []int {
	// 哈希表
	m := make(map[int]int)
	for i, v := range nums {
		if j, ok := m[target-v]; ok {
			return []int{j, i}
		}
		m[v] = i
	}
	return nil
}
