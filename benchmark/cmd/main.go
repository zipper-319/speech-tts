package main

import (
	"flag"
	"log"
	"speech-tts/benchmark"
	"sync"
)

var threadNum int
var useCaseNum int

func init() {
	flag.IntVar(&threadNum, "t", 1, "thread number, eg: -t 1")
	flag.IntVar(&useCaseNum, "u", 10, "useCase number, eg: -u 10")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	flag.Parse()
	log.Printf("thread number:%d; useCase number:%d", threadNum, useCaseNum)
	addr := "127.0.0.1:9000"
	text := "成都今天的天气"
	speaker := "DaXiaoFang"
	wg := sync.WaitGroup{}
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < useCaseNum; j++ {

				if err := benchmark.TestTTSV1(addr, text, speaker); err != nil {
					log.Println("_________")
					log.Printf("goroutine id:%d; err:%v", i, err)
					log.Println("_________")
					panic(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
