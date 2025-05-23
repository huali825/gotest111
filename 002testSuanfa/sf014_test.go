package _02testSuanfa

import (
	"fmt"
	"testing"
	"time"
)

//func Test014Name(t *testing.T) {
//	c := make(chan bool, 100)
//	for i := 0; i < 101; i++ {
//		go func(i int) {
//			fmt.Print(i, " ")
//			c <- true
//		}(i)
//	}
//
//	for i := 0; i < 100; i++ {
//		<-c
//	}
//}

func Test014Name(t *testing.T) {
	//fmt.Println("hello world")
	//defer func() {
	//	fmt.Println("defer hello world")
	//	if err := recover(); err != nil {
	//		fmt.Println(err)
	//	}
	//}()
	//panic("a panic is triggered")

	ch1 := make(chan int, 10)
	//ctx1 := context.Background()
	go func() {
		for i := 0; i < 100; i++ {
			ch1 <- i
			time.Sleep(time.Second)
		}
	}()

	for {
		time.Sleep(time.Second)
		v := <-ch1
		fmt.Println(v)
	}

}
