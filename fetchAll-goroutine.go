package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	startTime := time.Now()
	//创建string通道
	ch := make(chan string)
	//遍历参数传入的url
	for _, url := range os.Args[1:] {
		//启动一个goroutine
		go fetch(url, ch)
	}
	//从通道读数据
	for range os.Args[1:] {
		fmt.Println(<-ch)
	}
	//统计用时
	pass := time.Since(startTime).Seconds()
	fmt.Println("总用时%.2f秒", pass)
}

func fetch(url string, ch chan<- string) {
	startTime := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		//写入到通道
		ch <- fmt.Sprint(err)
		return
	}
	//拷贝到目的地 返回拷贝的字节数 第一个参数是目的地 这里是一个“丢弃”目的地
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	//关闭资源
	resp.Body.Close()
	if nil != err {
		//写入通道
		ch <- fmt.Sprintf("while reading %s: err:%v", url, err)
		return
	}
	//统计用时
	pass := time.Since(startTime).Seconds()
	//写入通道
	ch <- fmt.Sprintf("用时%.2f秒 数据量%7d字节 url:%s", pass, nbytes, url)
}
