package _02testSuanfa

import (
	"fmt"
	"testing"
)

// ListNode 定义链表节点结构体
type ListNode struct {
	Val  int
	Next *ListNode
}

// mergeTwoLists 合并两个有序链表
func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
	// 创建一个哨兵节点，用于简化操作
	dummy := &ListNode{}
	current := dummy

	// 遍历两个链表，比较节点值，将较小的节点添加到新链表中
	for l1 != nil && l2 != nil {
		if l1.Val <= l2.Val {
			current.Next = l1
			l1 = l1.Next
		} else {
			current.Next = l2
			l2 = l2.Next
		}
		current = current.Next
	}

	// 如果其中一个链表遍历完毕，将另一个链表的剩余部分添加到新链表中
	if l1 != nil {
		current.Next = l1
	} else {
		current.Next = l2
	}

	return dummy.Next
}

func Test005Name(t *testing.T) {
	// 创建两个有序链表
	l1 := &ListNode{Val: 1, Next: &ListNode{Val: 3, Next: &ListNode{Val: 5}}}
	l2 := &ListNode{Val: 2, Next: &ListNode{Val: 4, Next: &ListNode{Val: 6}}}

	// 合并两个有序链表
	mergedList := mergeTwoLists(l1, l2)

	// 打印合并后的链表
	for mergedList != nil {
		fmt.Printf("%d ", mergedList.Val)
		mergedList = mergedList.Next
	}
}
