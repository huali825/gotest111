package cronjob

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"testing"
	"time"
)

// 测试CronPractice001函数
func TestCronPractice001(t *testing.T) {
	// 创建一个新的Cron实例，并设置秒级精度
	c := cron.New(cron.WithSeconds())

	c.AddJob("@every 1s", MyJob{})

	// 添加一个定时任务，每隔7秒执行一次，打印当前时间
	_, err := c.AddFunc("0/4 * * * * *", func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	})
	// 如果添加任务失败，则返回
	if err != nil {
		return
	}

	// 启动Cron实例
	c.Start()
	defer c.Stop()

	// 等待300秒
	time.Sleep(300 * time.Second)

}

type MyJob struct {
}

func (j MyJob) Run() {
	log.Print("hello 开始运行了")
}
