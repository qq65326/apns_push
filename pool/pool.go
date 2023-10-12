/**
 * 任务池(生产者/消费者)
**/

package pool

import (
	"sync"
)

const MaxChanNum = 100000

var PoolChan = make(chan string, MaxChanNum)

// 设置数据
func SetOneNotification(notis string) {

	var wgp sync.WaitGroup // 声明一个信号量

	for i := 0; i < 1; i++ {
		wgp.Add(1) // 信号量加一
		go producer(&wgp, notis)
	}

	wgp.Wait() // 等待生产者退出，信号量为正时阻塞，直到信号量为0时被唤醒
}

// 入通道
func producer(wg *sync.WaitGroup, notis string) {

	PoolChan <- notis // 往通道里面放
	wg.Done()         // 信号量减一
}

// 获取数据
func GetOneNotification() string {

	var wgc sync.WaitGroup
	var notis string

	wgc.Add(1)
	go func() {
		notis = consumer(&wgc, PoolChan)
	}()

	wgc.Wait() // 等待消费者退出
	return notis
}

// 出通道
func consumer(wg *sync.WaitGroup, PoolChan chan string) string {

	var tmpNotis string

	product := <-PoolChan // // 从通道里面取
	tmpNotis = string(product)
	wg.Done()

	return tmpNotis
}
