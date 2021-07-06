package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
)

/**
创建一个downloader, 检查资源是否支持并发下载
 */

type Downloader struct {
	currCyn int
}

// 初始化一个downloader
func NewDownloader(currCyn int) *Downloader  {
	return &Downloader{
		currCyn: currCyn,
	}
}

// 增加一个Download方法, 首先发送一个head请求, 检查是否支持多线程下载
func (d *Downloader) Download (url, filename string) error  {
	if filename == "" {
		filename = path.Base(url)
	}
	res, err := http.Head(url)
	if err != nil {
		return err
	}
	fmt.Printf("head请求检查结果Accept-Ranges : %#v \n\n", res)
	// 如果http返回200的状态码, 并且 Accep-ranges 的值为bytes, 代表支持多线程下载
	if res.StatusCode == http.StatusOK && res.Header.Get("Accept-Ranges") == "bytes" {
		fmt.Printf("可以支持多线程下载\n")
		return d.multiDownload(url, filename, int(res.ContentLength))
	}
	fmt.Printf("不支持多线程下载\n")
	return d.singleDownload(url, filename)
}

// 多线程下载
func (d *Downloader) multiDownload(url, filename string, contentLength int) error {
	/*
		思路就是把大文件切片, 开启多个goroutine来并发下载切片文件, 然后合并文件
	*/

	// 1. 根据设定的线程数, 和文件的长度, 来计算每个线程承担多少文件长度的下载
	partSize := contentLength / d.currCyn

	// 2. 创建分片文件的存放目录
	os.Mkdir("./temp", 0777)
	//defer os.RemoveAll("./temp")

	var wg sync.WaitGroup // 声明一个等待组
	wg.Add(d.currCyn) // 为该等待组添加设定的几个任务
	rangeStart := 0
	for i := 0; i <= d.currCyn; i ++ {
		go func(i, rangeStart int) { // 匿名自执行函数
			defer wg.Done() // 每完成一个任务, 等待组里的任务-1

			// 每一个分片的长度
			rangeEnd := rangeStart + partSize

			// 最后一个分片, 总长度不能超过总文件的contentLength
			if i == d.currCyn-1 {
				rangeEnd = contentLength
			}
			// 开始分片下载
			d.downloadByPart(url, filename, rangeStart, rangeEnd, i)
		}(i, rangeStart)
		rangeStart += partSize + 1
	}

	// 等待所有的任务完成
	wg.Wait()


	// 合并下载的分片文件
	d.mergeFile(filename)

	return nil
}

func (d *Downloader) mergeFile(filename string) error{
	return nil
}

func (d *Downloader) downloadByPart(url, filename string, rangeStart, rangeEnd, i int) error  {

	if rangeStart >= rangeEnd {
		return nil
	}
	// 初始化一个http请求
	req, err := http.NewRequest("GET", url, nil)
	response, err := http.DefaultClient.Do(req)
	fmt.Printf("当前分片: %d, 内容: %#v \n\n", i, response)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	flags := os.O_CREATE | os.O_WRONLY // 这两个常量代表, 如果目录不存在就创建; O_WRONLY 采用只读的方式
	partFile, err := os.OpenFile(d.getPartFilename(filename, i), flags, 0666)

	fmt.Printf("partFile: %#v \n\n", partFile)
	if err != nil {
		log.Fatal(err)
	}
	defer partFile.Close()

	buf := make([]byte, 32*1024)
	//fmt.Printf("buf ====> %#v \n", buf)
	_, err = io.CopyBuffer(partFile, response.Body, buf)

	if err != nil {
		if err == io.EOF {
			return nil
		}
		log.Fatal(err)
	}
	return nil
}

// 单线程下载
func (d *Downloader) singleDownload(url, filename string) error  {
	return nil
}

func (d *Downloader) getPartFilename(filename string, index int) string {
	return fmt.Sprintf("./temp/%s-%d", filename, index)
}
