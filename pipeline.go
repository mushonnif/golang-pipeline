package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type FileInfo struct {
	FileName string
	Content  []byte
	Sum      string
	IsRename bool
}

func processFiles() {
	//pipeline 1: read all files
	chanFileContent := readFiles()

	//pipeline 2: calculate sum
	chanFileSum1 := getSum(chanFileContent)
	chanFileSum2 := getSum(chanFileContent)
	chanFileSum3 := getSum(chanFileContent)
	chanFileSum := mergeChanFileInfo(chanFileSum1, chanFileSum2, chanFileSum3)

	//pipeline 3: rename
	chanRename1 := rename(chanFileSum)
	chanRename2 := rename(chanFileSum)
	chanRename3 := rename(chanFileSum)
	chanRename4 := rename(chanFileSum)
	chanRename := mergeChanFileInfo(chanRename1, chanRename2, chanRename3, chanRename4)

	// print output
	counterRenamed := 0
	counterTotal := 0
	for fileInfo := range chanRename {
		if fileInfo.IsRename {
			counterRenamed++
		}
		counterTotal++
	}

	log.Printf("%d/%d files renamed", counterRenamed, counterTotal)
}

func readFiles() <-chan FileInfo {
	chanOut := make(chan FileInfo)
	go func() {
		err := filepath.Walk(tempPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			buf, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			chanOut <- FileInfo{
				FileName: path,
				Content:  buf,
			}

			return nil
		})
		if err != nil {
			fmt.Println("Error reading files")
		}

		close(chanOut)
	}()
	return chanOut
}

func getSum(chanIn <-chan FileInfo) <-chan FileInfo {
	chanOut := make(chan FileInfo)

	go func() {
		for fileInfo := range chanIn {
			sum := fmt.Sprintf("%x", md5.Sum(fileInfo.Content))
			fileInfo.Sum = sum

			chanOut <- fileInfo
		}
		close(chanOut)
	}()
	return chanOut
}

func mergeChanFileInfo(chanInMany ...<-chan FileInfo) <-chan FileInfo {
	wg := new(sync.WaitGroup)
	chanOut := make(chan FileInfo)

	wg.Add(len(chanInMany))
	for _, eachChan := range chanInMany {
		go func(chanIn <-chan FileInfo) {
			for fileInfo := range chanIn {
				chanOut <- fileInfo
			}
			wg.Done()
		}(eachChan)
	}

	go func() {
		wg.Wait()
		close(chanOut)
	}()

	return chanOut
}

func rename(chanIn <-chan FileInfo) <-chan FileInfo {
	chanOut := make(chan FileInfo)

	go func() {
		for fileInfo := range chanIn {
			newPath := filepath.Join(tempPath, fmt.Sprintf("file-%s.txt", fileInfo.Sum))
			err := os.Rename(fileInfo.FileName, newPath)
			fileInfo.IsRename = err == nil

			chanOut <- fileInfo
		}

		close(chanOut)
	}()
	return chanOut
}
