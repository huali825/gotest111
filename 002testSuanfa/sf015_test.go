/*
date:2025年4月21日21:00:27
title:cond的使用
author:tmh
*/

package _02testSuanfa

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Test015Name 测试队列
// 测试函数，用于测试队列的入队和出队操作
func Test015Name(t *testing.T) {
	// 创建一个队列
	q := &Queue{
		cond:  sync.NewCond(&sync.Mutex{}), // 创建一个互斥锁
		queue: []string{},                  // 创建一个空的队列
	}
	// 启动一个goroutine，用于向队列中添加元素
	go func() {
		for {
			q.Enqueue("a")              // 向队列中添加元素
			time.Sleep(time.Second * 2) // 每隔2秒添加一个元素
		}
	}()
	// 无限循环，用于从队列中取出元素
	for {
		q.Dequeue()                 // 从队列中取出元素
		time.Sleep(time.Second * 1) // 每隔1秒取出一个元素
	}
}

type Queue struct {
	queue []string
	cond  *sync.Cond // 条件变量
}

// Enqueue 入队
func (q *Queue) Enqueue(s string) {
	q.cond.L.Lock()                // 加锁
	defer q.cond.L.Unlock()        // 解锁
	q.queue = append(q.queue, s)   // 将元素添加到队列中
	fmt.Printf("Enqueue: %s\n", s) // 打印入队元素
	q.cond.Broadcast()             // 通知所有等待的goroutine
}

// Dequeue 出队
func (q *Queue) Dequeue() string {
	q.cond.L.Lock()         // 加锁
	defer q.cond.L.Unlock() // 解锁
	if len(q.queue) == 0 {  // 如果队列为空
		fmt.Println("Queue is empty, waiting...") // 打印提示信息
		q.cond.Wait()                             // 等待直到有新的元素入队
	}
	result := q.queue[0]  // 取出队列中的第一个元素
	q.queue = q.queue[1:] // 将队列中的第一个元素移除
	return result         // 返回取出的元素
}
