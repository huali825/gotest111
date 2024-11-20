package _02testSuanfa

import (
	"fmt"
	"testing"
)

func Test001Name(t *testing.T) {
	//10 100 11 101

	temp := funcTest001Name01()
	fmt.Println(temp)

	temp2 := funcTest001Name02()
	fmt.Println(temp2)

	temp3 := funcTest001Name03()
	fmt.Println(temp3)

	temp4 := funcTest001Name04()
	fmt.Println(temp4)
}

func funcTest001Name01() int {
	val := 10
	defer func() {
		val += 1
	}()
	return val
}

func funcTest001Name02() int {
	val := 10
	defer func() {
		val += 1
	}()
	return 100
}

func funcTest001Name03() (val int) {
	val = 10
	defer func() {
		val += 1
	}()
	return val
}

func funcTest001Name04() (val int) {
	val = 10
	defer func() {
		val += 1
	}()
	return 100
}
