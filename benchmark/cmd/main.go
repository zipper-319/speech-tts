package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"speech-tts/benchmark"
	"speech-tts/internal/utils"
	"strings"
	"sync"
	"sync/atomic"
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
var filePath string
var outfile string

func init() {
	flag.IntVar(&threadNum, "t", 1, "thread number, eg: -t 1")
	flag.IntVar(&useCaseNum, "u", 10, "useCase number, eg: -u 10")
	flag.StringVar(&addr, "a", "127.0.0.1:3012", "addr, eg: -a 127.0.0.1:3012")
	flag.StringVar(&speaker, "s", "DaXiaoFang", "speaker name, eg: -s DaXiaoFang")
	flag.StringVar(&testVersion, "v", "v2", "test Version, eg: -v v1")
	flag.StringVar(&movement, "m", "Nvidia-a2g", "movement, eg: -m SweetGirl")
	flag.StringVar(&expression, "e", "", "expression, eg: -e FaceGood")
	flag.BoolVar(&isSaveFile, "i", false, "isSaveFile, eg: -i true")
	flag.StringVar(&filePath, "f", "./testText.txt", "filePath, eg: -f ./testText.txt")
	flag.StringVar(&outfile, "o", "./testResult.txt", "outfile, eg: -o ./testResult.txt")
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Flags())
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	flag.Parse()
	log.Printf("thread number:%d; useCase number:%d, speaker:%s, testVersion:%s, movement:%s", threadNum, useCaseNum, speaker, testVersion, movement)
	v2version := benchmark.GetV2Version(context.Background(), addr)
	outResultList := atomic.Value{}
	outResultList.Store(make([]*benchmark.OutResult, 0, useCaseNum))

	out, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer func() {
		resultList := outResultList.Load().([]*benchmark.OutResult)
		var totalFirstClientCost, totalClientCost, totalFirstServerCost, totalServerCost int64
		for _, outResult := range resultList {
			totalFirstClientCost += outResult.FirstClientCost
			totalClientCost += outResult.ClientCost
			totalFirstServerCost += int64(outResult.FirstServerCost)
			totalServerCost += int64(outResult.ServerCost)
		}
		averageFirstClientCost := totalFirstClientCost / int64(len(resultList))
		averageClientCost := totalClientCost / int64(len(resultList))
		averageFirstServerCost := totalFirstServerCost / int64(len(resultList))
		averageServerCost := totalServerCost / int64(len(resultList))
		out.WriteString(fmt.Sprintf("\n 版本：%s \n", v2version))
		out.WriteString(fmt.Sprintf(" 用例总数：%d  并发数：%d 发音人：%s tts服务端地址：%s 客户端第一帧平均耗时：%d, 客户端平均耗时：%d, 服务端第一帧平均耗时:%d 服务端平均耗时:%d \n", useCaseNum, threadNum, speaker, addr, averageFirstClientCost, averageClientCost, averageFirstServerCost, averageServerCost))
		out.Close()
	}()

	out.WriteString(fmt.Sprintf("%s  %s  %s  %s  %s\n", "文本", "客户端第一帧耗时", "客户端总耗时", "服务端第一帧耗时", "服务端总耗时"))

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
					if outResult, err := benchmark.TestTTSV2(ctx, out, user, addr, text, speaker, fmt.Sprintf("test_thread%d_%dnum", t, num), fmt.Sprintf("test_robot_thread%d_%dnum", t, num),
						movement, expression, num, isSaveFile); err == nil {
						resultList := outResultList.Load().([]*benchmark.OutResult)
						out.WriteString(fmt.Sprintf("%s   %dms   %dms     %dms     %dms\n", outResult.Text, outResult.FirstClientCost, outResult.ClientCost, outResult.FirstServerCost, outResult.ServerCost))
						resultList = append(resultList, outResult)
						outResultList.Store(resultList)
					} else {
						log.Println("_________")
						log.Printf("TestTTSV2; goroutine id:%d; err:%v", i, err)
						log.Println("_________")
					}
				}
			}
		}(i)

	}
	var currentLine int

	for {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("打开文件时出错: %v\n", err)
			return
		}

		scanner := bufio.NewScanner(file)
		// 循环读取每一行直到EOF或达到指定行数
		for scanner.Scan() && currentLine < useCaseNum {
			ch <- scanner.Text()
			currentLine++
		}

		if err = scanner.Err(); err != nil {
			fmt.Printf("读取文件时出错: %v\n", err)
			file.Close()
			return
		}

		file.Close()

		if currentLine >= useCaseNum {
			break // 如果已经达到指定的行数，则退出循环
		}
	}
	close(ch)
	wg.Wait()
}
