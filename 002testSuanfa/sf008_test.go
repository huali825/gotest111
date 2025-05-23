package _02testSuanfa

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test008Name(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	hchan := make(chan int, 10)
	var wg sync.WaitGroup // 创建一个WaitGroup，用于等待所有goroutine完成
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(ctx, &wg, hchan)
	}

	for i := 1; i <= 10; i++ {
		hchan <- i
	}
	close(hchan)
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("All tasks processed successfully") // 所有任务都成功处理
	case <-ctx.Done():
		fmt.Println("Context timeout or cancellation") // 超时或取消
	}
}

func worker(ctx context.Context, wg *sync.WaitGroup, tasks <-chan int) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case taskID, ok := <-tasks: // 从任务通道中接收任务ID，ok表示通道是否已关闭
			if !ok {
				return
			}
			time.Sleep(time.Duration(taskID) * 100 * time.Millisecond)
			fmt.Printf("Processed task %d\n", taskID)
		}
	}
}
