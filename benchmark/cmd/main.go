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
var speaker string
var testVersion string
var movement string
var expression string

func init() {
	flag.IntVar(&threadNum, "t", 1, "thread number, eg: -t 1")
	flag.IntVar(&useCaseNum, "u", 10, "useCase number, eg: -u 10")
	flag.StringVar(&addr, "a", "127.0.0.1:3012", "addr, eg: -a 127.0.0.1:3012")
	flag.StringVar(&speaker, "s", "DaXiaoFang", "speaker name, eg: -s DaXiaoFang")
	flag.StringVar(&testVersion, "v", "v1", "test Version, eg: -v v1")
	flag.StringVar(&movement, "m", "", "movement, eg: -m SweetGirl")
	flag.StringVar(&expression, "e", "", "expression, eg: -e FaceGood")
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Flags())
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	flag.Parse()
	log.Printf("thread number:%d; useCase number:%d, speaker:%s, testVersion:%s", threadNum, useCaseNum, speaker, testVersion)

	text := "成都今天的天气"
	wg := sync.WaitGroup{}
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go func(t int) {
			for j := 0; j < useCaseNum; j++ {

				if testVersion == "v1" {
					if err := benchmark.TestTTSV1(addr, text, speaker, fmt.Sprintf("test_thread%d_%d", t, j), fmt.Sprintf(fmt.Sprintf("test_robot_thread%d_%d", t, j)), j); err != nil {
						log.Println("_________")
						log.Printf("goroutine id:%d; err:%v", i, err)
						log.Println("_________")
						panic(err)
					}
				}

				if testVersion == "v2" {
					if err := benchmark.TestTTSV2(addr, text, speaker, fmt.Sprintf("test_thread%d_%d", t, j), fmt.Sprintf(fmt.Sprintf("test_robot_thread%d_%d", t, j)),
						movement, expression, j); err != nil {
						log.Println("_________")
						log.Printf("goroutine id:%d; err:%v", i, err)
						log.Println("_________")
						panic(err)
					}
				}

			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
