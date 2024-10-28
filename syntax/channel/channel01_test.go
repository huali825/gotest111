package channel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	fmt.Println("hello test_channel01  " + time.Now().String())
}

// 定义一个结构体MyStruct，包含一个通道ch和一个sync.Once类型的closeOnce
type MyStruct struct {
	ch        chan struct{}
	closeOnce sync.Once
}

// 定义一个Close方法，用于关闭通道ch
func (s *MyStruct) Close() error {
	// 使用sync.Once类型的closeOnce来确保通道ch只被关闭一次
	s.closeOnce.Do(func() {
		close(s.ch)
	})
	// 返回nil，表示关闭成功
	return nil
}

func TestForLoop(t *testing.T) {
	// 创建一个长度为10的int类型切片
	slice := make([]int, 10, 44)
	fmt.Println(slice)
	fmt.Println(&slice)

	// 创建一个初始大小为10的映射
	map1 := make(map[string]int)
	fmt.Println(map1)
	fmt.Println(&map1)

	// 创建一个容量为10的缓冲通道
	ch := make(chan int, 10)
	fmt.Println(ch)
	fmt.Println(&ch)

	// 使用 new() 创建一个 int 类型的零值变量的指针
	numPtr := new(int)
	fmt.Println(*numPtr) // 输出 0	（int 类型的零值）
	fmt.Println(numPtr)  // 输出 0xc00008c308 （一个内存地址）
}
