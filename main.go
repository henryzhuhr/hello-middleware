package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, Database!")

	// 获取当前时间戳
	now := time.Now()
	fmt.Println("Now:", now)

	// 时间戳转换为字符串
	var tm time.Time
	// tm = time.Unix(now.Unix(), 0)
	tm = time.Unix(1728386092, 0)
	fmt.Println("Time:", tm)

}
