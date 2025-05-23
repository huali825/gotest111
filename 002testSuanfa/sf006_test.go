package _02testSuanfa

import (
	"testing"
)

func Test006Name(t *testing.T) {
	// 测试代码
	nums := []int{3, 2, 2, 3, 2, 2, 2, 2}
	t.Log(removeElement(nums, 3))
	t.Log(nums)

	nums = []int{1, 2, 3, 1, 1, 2}
	t.Log(numIdenticalPairs(nums))

}

func removeElement(nums []int, val int) int {
	k := 0
	for _, x := range nums {
		if x != val {
			nums[k] = x
			k++
		}
	}
	return k
}

// numIdenticalPairs 函数用于计算数组中相同元素的数对数量
func numIdenticalPairs(nums []int) (ans int) {
	// cnt 是一个映射，用于记录每个数字出现的次数
	cnt := map[int]int{}

	// 遍历数组中的每个元素
	for _, x := range nums {
		// 将当前元素 x 在 cnt 中出现的次数加到 ans 中
		// 这是因为如果 x 出现了 n 次，那么它可以形成 n*(n-1)/2 个数对
		ans += cnt[x]

		// 将当前元素 x 的计数加 1
		cnt[x]++
	}

	// 返回计算出的数对数量
	return
}
