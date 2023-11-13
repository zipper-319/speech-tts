package util

import (
	"github.com/go-kratos/kratos/v2/log"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func DownloadFile(url, fileName string) error {
	log.Debug("Starting download...")
	wg := sync.WaitGroup{}
	parts := 4

	// 获取文件头信息
	head, err := http.Head(url)
	if err != nil {
		return err
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

	log.Debug("Download complete!")
	return nil
}

func Download(url string, start int64, end int64, index int, fileName string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err)
		return
	}

	// 指定Range头
	rangeHeader := "bytes=" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10)
	request.Header.Set("Range", rangeHeader)

	resp, err := client.Do(request)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Error(err)
		return
	}
	defer file.Close()

	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		log.Error(err)
		return
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("Part %d finished\n", index)
}
