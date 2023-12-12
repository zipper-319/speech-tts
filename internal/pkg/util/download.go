package util

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func DownloadChunk(url, fileName string) error {
	log.Printf("Starting download...; url:%s; fileName:%s", url, fileName)
	wg := sync.WaitGroup{}
	parts := 4

	// 获取文件头信息
	head, err := http.Head(url)
	if err != nil {
		return err
	}
	if head.StatusCode != 200 {
		return errors.New("http status code is not 200; statusCode:" + strconv.Itoa(head.StatusCode))
	}
	filesize := head.ContentLength

	// 创建文件
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	partSize := filesize / int64(parts)

	for i := 0; i < parts; i++ {
		startIndex := int64(i) * partSize
		endIndex := startIndex + partSize - 1
		if i == parts-1 {
			endIndex = filesize - 1
		}
		wg.Add(1)
		go Download(url, startIndex, endIndex, i, fileName, &wg)
	}
	wg.Wait()

	log.Println("Download complete!")
	return nil
}

func Download(url string, start int64, end int64, index int, fileName string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 指定Range头
	rangeHeader := "bytes=" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10)
	request.Header.Set("Range", rangeHeader)

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Part %d finished\n", index)
}

func DownloadFile(url, filepath string) error {
	// 发起GET请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建一个空文件用于保存下载的数据
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 将响应流直接复制到文件中
	_, err = io.Copy(out, resp.Body)
	return err // 如果成功则err为nil，否则返回错误信息
}
