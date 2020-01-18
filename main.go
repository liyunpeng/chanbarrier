package main

import (
	"fmt"
	"time"
)

type Barrier interface {
	Wait ()
}

type barrier struct {
	chCount chan struct{} // 所有调用Wait()函数的先通过这个chan阻塞
	count   int           // 记住数量，阻塞超过这个量就可以激活所有协程了
	chSync  chan struct{}     // 增加一个chan
}
func NewBarrier(n int) Barrier {
	b := &barrier{count: n, chCount: make(chan struct{}), chSync: make(chan struct{})}
	go b.Sync()
	return b
}
func (b *barrier) Wait() {
	b.chCount <- struct{}{}
	fmt.Println("b.chCount <- struct{}{}")
	<-b.chSync                // 再次阻塞
}
func (b *barrier) Sync() {
	count := 0
	fmt.Println("(b *barrier) Sync()")
	for range b.chCount {
		count++
		if count >= b.count {
			fmt.Println("close %v", b.chSync)
			close(b.chSync)   // close这个chan所有阻塞协程都会被激活
			break
		}
	}
}

// 测试代码
func main () {
	// 创建栅栏对象
	b := NewBarrier(10)
	// 达到的效果：前9个协程调用Wait()阻塞，第10个调用后10个协程全部唤醒
	for i:=0; i<10; i++ {
		fmt.Println(i)
		go b.Wait()
	}

	time.Sleep(100000)
}