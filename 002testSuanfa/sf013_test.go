package _02testSuanfa

import (
	"fmt"
	"testing"
)

// 2025年4月15日02:13:13
func Test013Name(t *testing.T) {
	mylist := []string{"i", "am", "stupid", "and", "weak"}
	fmt.Println(mylist)
	for k, v := range mylist {
		if v == "stupid" {
			mylist[k] = "smart"
		}
		if v == "weak" {
			mylist[k] = "strong"
		}
	}
	fmt.Println(mylist)

}
