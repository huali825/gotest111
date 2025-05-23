package main

import (
	"fmt"
	"time"
)

func main() {

	for i := 0; i < 100; i++ {
		go func(temp int, tempAdd int) {
			fmt.Println(temp, tempAdd)
		}(i, i+1)
	}

	time.Sleep(time.Second)
}
