package util

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func DeCompress(zipFileName string) error {
	// 第一步，打开 zip 文件
	zipFile, err := zip.OpenReader(zipFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 第二步，遍历 zip 中的文件
	for _, f := range zipFile.File {
		fileName := f.Name
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(fileName, os.ModePerm)
			continue
		}
		// 创建对应文件夹
		if err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
			return err
		}
		// 解压到的目标文件
		dstFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		file, err := f.Open()
		if err != nil {
			return err
		}
		// 写入到解压到的目标文件
		if _, err := io.Copy(dstFile, file); err != nil {
			return err
		}
		dstFile.Close()
		file.Close()
	}
	return nil
}

func DeCompressToPath(zipFileName, dstPath string) error {
	// 第一步，打开 zip 文件
	zipFile, err := zip.OpenReader(zipFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
		return err
	}

	// 第二步，遍历 zip 中的文件
	for _, f := range zipFile.File {

		if f.FileInfo().IsDir() {
			continue
		}
		fileName := strings.Split(f.Name, "/")[len(strings.Split(f.Name, "/"))-1]
		if fileName == "" {
			return errors.New("file name is empty")
		}
		fileName = dstPath + "/" + fileName
		// 创建对应文件夹
		if err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
			return err
		}
		// 解压到的目标文件
		dstFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		file, err := f.Open()
		if err != nil {
			return err
		}
		// 写入到解压到的目标文件
		if _, err := io.Copy(dstFile, file); err != nil {
			return err
		}
		dstFile.Close()
		file.Close()
	}
	return nil
}
