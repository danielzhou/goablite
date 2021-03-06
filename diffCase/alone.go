/*
# 1. 单机版，args取参数

参考：https://www.jianshu.com/p/7f28aad6ecc5

@brief:只能输入2个参数：并发数 目标url

介绍：输入2个参数：并发数 目标url

使用方法：在对应目录下，使用go build进行编译，然后运行alone.exe n url，不限连接次数，需手动停止。

*/
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	SBCNum     int           // 并发连接数
	QPSNum     int           // 总请求次数
	RTNum      time.Duration // 响应时间
	RTTNum     time.Duration // 平均响应时间
	SuccessNum int           // 成功次数
	FailNum    int           // 失败次数

	BeginTime time.Time // 开始时间
	SecNum    int       // 秒数

	RQNum int    // 最大并发数，由命令行传入
	Url   string // url，由命令行传入

	controlNum chan int // 控制并发数量
)

var mu sync.Mutex // 加锁

func init() {
	if len(os.Args) != 3 {
		fmt.Println(len(os.Args))
		fmt.Println(os.Args)
		log.Fatal("请求次数 url")
	}
	RQNum, _ = strconv.Atoi(os.Args[1])
	controlNum = make(chan int, RQNum)
	Url = os.Args[2]
}

func main() {
	go func() {
		for range time.Tick(1 * time.Second) {
			SecNum++
			fmt.Printf("并发数：%d，请求次数：%d,平均响应时间：%s,成功次数：%d,失败次数：%d\n",
				len(controlNum), SuccessNum+FailNum, RTNum/(time.Duration(SecNum)*time.Second),
				SuccessNum, FailNum)
		}
	}()
	requite()
}

func requite() {
	for {
		controlNum <- 1
		go func(c chan int) {
			var tb time.Time
			var el time.Duration
			for {
				tb = time.Now()
				_, err := http.Get(Url)
				if err == nil {
					el = time.Since(tb)
					mu.Lock() // 上锁
					SuccessNum++
					RTNum += el
					mu.Unlock() // 解锁
				} else {
					mu.Lock() // 上锁
					FailNum++
					mu.Unlock() // 解锁
				}
				time.Sleep(1 * time.Second)
			}
			<-c
		}(controlNum)
		time.Sleep(45 * time.Millisecond)
	}
}
