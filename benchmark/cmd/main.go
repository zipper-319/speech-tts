package main

import (
	"log"
	"speech-tts/benchmark"
	"sync"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	addr := "127.0.0.1:9000"
	text := "成都今天的天气"
	speaker := "DaXiaoFang"
	wg := sync.WaitGroup{}
	for i := 0; i <= 5; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j <= 10000; j++ {

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
