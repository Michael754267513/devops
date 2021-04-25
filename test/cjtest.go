package main

import (
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func main() {
	log.Println("Starting...")

	// 定义一个cron运行器
	c := cron.New()
	// 定时5秒，每5秒执行print5
	c.AddFunc("*/5 * * * * *", print5)
	// 定时15秒，每5秒执行print5
	c.AddFunc("*/15 * * * * *", print15)

	// 开始
	c.Start()
	defer c.Stop()

	// 这是一个使用time包实现的定时器，与cron做对比
	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
			print10()
		}
	}
}

func print5() {
	log.Println("Run 5s cron")
}

func print10() {
	log.Println("Run 10s cron")
}

func print15() {
	log.Println("Run 15s cron")
}
