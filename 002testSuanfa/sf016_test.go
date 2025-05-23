/*
date:2025年4月22日09:35:59
title: sync.once
author:tmh
*/

package _02testSuanfa

import (
	"sync"
	"testing"
)

func Test016Name(t *testing.T) {
	o := sync.Once{}
	slice1 := []int{1, 2, 3, 4, 5}
	o.Do(func() {
		slice1 = add(slice1, 11)
	})
	o.Do(func() {
		slice1 = add(slice1, 12)
	})
	o.Do(func() {
		slice1 = add(slice1, 13)
	})
	t.Log(slice1)
}

func add(a []int, b int) []int {
	return append(a, b)
}
