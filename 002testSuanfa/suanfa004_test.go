package _02testSuanfa

import (
	"fmt"
	"testing"
)

func TestSuanfa004(t *testing.T) {
	var temp, newlen, oldcap = 0, 0, 0

	t.Log("原切片小于256,需求切片 大于 两倍原切片")
	temp, newlen, oldcap = 0, 220, 100
	temp = nextslicecap(newlen, oldcap)
	fmt.Println("原切片,需求切片,实际算的切片", oldcap, newlen, temp)

	t.Log("原切片大于256,需求切片 大于 两倍原切片")
	temp, newlen, oldcap = 0, 2200, 1000
	temp = nextslicecap(newlen, oldcap)
	fmt.Println("原切片,需求切片,实际算的切片", oldcap, newlen, temp)

	t.Log("原切片小于256,需求切片 小于 两倍原切片")
	temp, newlen, oldcap = 0, 9, 8
	temp = nextslicecap(newlen, oldcap)
	fmt.Println("原切片,需求切片,实际算的切片", oldcap, newlen, temp)

	t.Log("原切片大于256,需求切片 小于 两倍原切片")
	temp, newlen, oldcap = 0, 1100, 1000
	temp = nextslicecap(newlen, oldcap)
	fmt.Println("原切片,需求切片,实际算的切片", oldcap, newlen, temp)

	t.Log("原切片大于256,需求切片 小于 两倍原切片")
	temp, newlen, oldcap = 0, 11000, 10000
	temp = nextslicecap(newlen, oldcap)
	fmt.Println("原切片,需求切片,实际算的切片", oldcap, newlen, temp)

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

// nextslicecap函数用于计算切片的下一个容量
func nextslicecap(newLen, oldCap int) int {
	// 将newcap初始化为oldcap
	newcap := oldCap
	// 计算doublecap，即oldcap的两倍
	doublecap := newcap + newcap
	// 如果newLen大于doublecap，则返回newLen
	if newLen > doublecap {
		return newLen
	}

	// 定义阈值，当oldcap小于阈值时，使用doublecap作为newcap
	const threshold = 256
	if oldCap < threshold {
		return doublecap
	}
	// 循环计算newcap
	for {
		// 当切片较小时，每次增长2倍
		// 当切片较大时，每次增长1.25倍+192
		// 该公式给出了平滑的过渡
		newcap += (newcap + 3*threshold) >> 2

		// 需要检查newcap >= newLen和newcap是否溢出
		// newLen保证大于零，因此当newcap溢出时，uint(newcap) > uint(newLen)
		// 这允许我们使用相同的比较来检查两者
		if uint(newcap) >= uint(newLen) {
			break
		}
	}

	// 当newcap计算溢出时，将newcap设置为请求的cap
	if newcap <= 0 {
		return newLen
	}
	return newcap
}
