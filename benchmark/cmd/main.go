package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"math/rand"
	"os"
	"speech-tts/benchmark"
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
			for {
				text, ok := <-ch
				if !ok {
					return
				}
				text = strings.TrimSpace(text)
				if text == "" {
					continue
				}
				num += 1
				md := metadata.Pairs(
					"authorization", "Bearer some-secret-token",
				)
				ctxBase := metadata.NewOutgoingContext(context.Background(), md)
				ctx, cancel := context.WithCancel(ctxBase)

				go func() {
					rand.Seed(time.Now().UnixNano())
					n := rand.Int31n(1000) + 100
					time.Sleep(time.Millisecond * time.Duration(n))
					cancel()
					log.Printf("cancel after %d ms ", n)
				}()

				if time.Now().Unix()%2 == 0 {
					if err := benchmark.TestTTSV1(ctx, addr, text, speaker, fmt.Sprintf("test_thread%d", t), fmt.Sprintf("test_robot_thread%d", t), num); err != nil {
						log.Println("_________")
						log.Printf("TestTTSV1; goroutine id:%d; err:%v", i, err)
						log.Println("_________")
						panic(err)
					}
				} else {
					if err := benchmark.TestTTSV2(ctx, addr, text, speaker, fmt.Sprintf("test_thread%d", t), fmt.Sprintf("test_robot_thread%d", t),
						movement, expression, num); err != nil {
						log.Println("_________")
						log.Printf("TestTTSV2; goroutine id:%d; err:%v", i, err)
						log.Println("_________")
						panic(err)
					}
				}
			}
		}(i)

	}
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
	}

	wg.Wait()
}
