package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"os"
	"speech-tts/benchmark"
	"speech-tts/internal/utils"
	"strings"
	"sync"
	"time"
)

var threadNum int
var useCaseNum int
var addr string
var speaker string
var testVersion string
var movement string
var expression string
var isSaveFile bool

func init() {
	flag.IntVar(&threadNum, "t", 1, "thread number, eg: -t 1")
	flag.IntVar(&useCaseNum, "u", 10, "useCase number, eg: -u 10")
	flag.StringVar(&addr, "a", "127.0.0.1:3012", "addr, eg: -a 127.0.0.1:3012")
	flag.StringVar(&speaker, "s", "DaXiaoFang", "speaker name, eg: -s DaXiaoFang")
	flag.StringVar(&testVersion, "v", "", "test Version, eg: -v v1")
	flag.StringVar(&movement, "m", "Nvidia-a2g", "movement, eg: -m SweetGirl")
	flag.StringVar(&expression, "e", "", "expression, eg: -e FaceGood")
	flag.BoolVar(&isSaveFile, "i", false, "isSaveFile, eg: -i true")
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Flags())

}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	flag.Parse()
	log.Printf("thread number:%d; useCase number:%d, speaker:%s, testVersion:%s, movement:%s", threadNum, useCaseNum, speaker, testVersion, movement)

	file, err := os.Open("./testTTS.txt")
	reader := bufio.NewReader(file)
	if err != nil {
		return
	}
	ch := make(chan string, 100)

	wg := sync.WaitGroup{}
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go func(t int) {
			defer wg.Done()
			num := 0
			currentTestVersion := "v1"
			for {
				text, ok := <-ch
				if !ok {
					return
				}
				text = strings.TrimSpace(text)
				if text == "" {
					continue
				}
				if testVersion == "" {

					if time.Now().Unix()%2 == 0 {
						currentTestVersion = "v1"
					} else {
						currentTestVersion = "v2"
					}
				} else {
					currentTestVersion = testVersion
				}
				num += 1
				md := metadata.Pairs(
					"Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lSWQiOjAsIkFjY291bnQiOiJIQVJJWF9BSV9TSE9XX1BMQVRGT1JNIiwiUm9sZSI6IiIsImF1ZCI6WyJIQVJJWF9BSV9TSE9XX1BMQVRGT1JNIl0sImV4cCI6MTg0NjU2NDcwMCwiaWF0IjoxNjkxMDQ0NzAwLCJqdGkiOiI1NWY0N2U5Ni0zMWM4LTExZWUtOGZmNi00YWFlYTEwZTg5MWMifQ.Aj0dadEH1aXSIvz6RWGtmFXSbzY-QQS_-9jEDpB4IYU",
				)
				ctxBase := metadata.NewOutgoingContext(context.Background(), md)
				ctx, _ := context.WithCancel(ctxBase)

				//go func() {
				//	rand.Seed(time.Now().UnixNano())
				//	n := rand.Int31n(1000) + 100
				//	time.Sleep(time.Millisecond * time.Duration(n))
				//	cancel()
				//	log.Printf("cancel after %d ms ", n)
				//}()

				if currentTestVersion == "v1" {

					if err := benchmark.TestTTSV1(ctx, addr, text, speaker, fmt.Sprintf("test_thread%d_%dnum", t, num), fmt.Sprintf("test_robot_thread%d_%dnum", t, num), num); err != nil {
						log.Println("_________")
						log.Printf("TestTTSV1; goroutine id:%d; err:%v", i, err)
						log.Println("_________")
					}

				} else {
					user := utils.DefaultUser
					if err := benchmark.TestTTSV2(ctx, user, addr, text, speaker, fmt.Sprintf("test_thread%d_%dnum", t, num), fmt.Sprintf("test_robot_thread%d_%dnum", t, num),
						movement, expression, num, isSaveFile); err != nil {
						log.Println("_________")
						log.Printf("TestTTSV2; goroutine id:%d; err:%v", i, err)
						log.Println("_________")
					}

				}
			}
		}(i)

	}
	var lineNum int
	for {
		line, _, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			log.Println(err)
			return
		}
		if err == io.EOF {
			break
		}
		ch <- string(line)
		lineNum += 1
		if lineNum == useCaseNum {
			log.Printf("finished to test tts;num:%d", useCaseNum)
			close(ch)
			break
		}
	}

	wg.Wait()
}
