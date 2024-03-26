package weixin

import (
	"fmt"
	"testing"
	"time"
)

func TestTicket(t *testing.T) {
	ticker := time.NewTicker(time.Second * 1)

	ticker.Reset(time.Second)
	i := 0
	for i < 10 { //循环
		i++
		fmt.Println("i = ", i)

		if i == 5 {
			ticker.Stop() //停止定时器
			break
		}
	}

	i = 0
	ticker.Reset(time.Second)

	for i < 10 { //循环
		<-ticker.C
		i++
		fmt.Println("i = ", i)

		if i == 5 {
			ticker.Stop() //停止定时器
			break
		}
	}
}
