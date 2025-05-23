/*
date:2025年4月22日11:13:59
title: 多个消费者多个生产者,大家都是间隔一秒发送/接收
author:tmh
*/

package _02testSuanfa

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test017Name(t *testing.T) {
	sc := &safeCh1{
		ch:     make(chan int, 10),
		outNum: 0,
		mutex:  sync.Mutex{},
	}

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(ProducerID int) {
			defer wg.Done()
			for {
				sc.mutex.Lock()
				value := sc.outNum
				sc.outNum++
				sc.mutex.Unlock()

				sc.ch <- value
				fmt.Printf("Producer %d: sent %d (Queue len: %d)\n", ProducerID, value, len(sc.ch))
				time.Sleep(time.Second)
			}
		}(i)
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			for {
				value := <-sc.ch
				fmt.Printf("Consumer %d: received %d (Queue len: %d)\n", consumerID, value, len(sc.ch))
				time.Sleep(time.Second)
			}
		}(i)
	}

	wg.Wait()
	close(sc.ch)
	fmt.Println("Program finished.")
}

type safeCh1 struct {
	ch     chan int
	mutex  sync.Mutex
	outNum int
}
