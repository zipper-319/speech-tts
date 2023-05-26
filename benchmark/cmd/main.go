package main

import (
	"flag"
	"fmt"
	"log"
	"speech-tts/benchmark"
	"sync"
)

var threadNum int
var useCaseNum int
var addr string

func init() {
	flag.IntVar(&threadNum, "t", 1, "thread number, eg: -t 1")
	flag.IntVar(&useCaseNum, "u", 10, "useCase number, eg: -u 10")
	flag.StringVar(&addr, "a", "127.0.0.1:3012", "addr, eg: -a 127.0.0.1:3012")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	flag.Parse()
	log.Printf("thread number:%d; useCase number:%d", threadNum, useCaseNum)

	text := "成都今天的天气"
	speaker := "DaXiaoFang"
	wg := sync.WaitGroup{}
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go func(t int) {
			for j := 0; j < useCaseNum; j++ {

				if err := benchmark.TestTTSV2(addr, text, speaker, fmt.Sprintf("test_thread%d_%d", t, j), fmt.Sprintf(fmt.Sprintf("test_robot_thread%d_%d", t, j))); err != nil {
					log.Println("_________")
					log.Printf("goroutine id:%d; err:%v", i, err)
					log.Println("_________")
					panic(err)
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
