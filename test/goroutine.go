package main

import (
	"fmt"
	"time"
)

func count(thing string) {
	for i := 0; i < 10; i++ {
		fmt.Println(i, thing)
		time.Sleep(time.Microsecond * 500)
	}
}

func count2(thing string) {
	for i := 0; i < 5; i++ {
		fmt.Println(i, thing)
		time.Sleep(time.Microsecond * 500)
	}
}
func testGoRoutine() {
	c := make(chan string)
	c <- "hello" // If the channel is unbuffered, the sender blocks until the receiver has received the value
	msg := <-c
	fmt.Println(msg)

}

func testChannel() {
	c1 := make(chan string)
	c2 := make(chan string)
	go func() {
		for {
			c1 <- "500ms"
			time.Sleep(time.Millisecond * 500)
		}
	}()

	go func() {
		for {
			c2 <- "3000ms"
			time.Sleep(time.Millisecond * 3000)
		}
	}()

	for {
		msg1 := <-c1
		fmt.Println(msg1)
		msg2 := <-c2
		fmt.Println(msg2)
	}
}

//1 立即打印出 500 ，3000
//2 随后打印出500 此时执行的是msg1 := <- c1 fmt.Println(msg1)和channel1里的代码
//3 当遇到msg2 := <- c2时， channel2还在sleep， 所以过3s后再打印3000
// 重复2-3
// 500不能被立即打印因为要等 channel2读写完毕才能进行新的一轮
