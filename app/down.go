package main

import (
	"bufio"
	"fmt"
	"net/http"
	"io"
	"os"
	"strings"
	"sync"
)

var (
	//下载目录
	downloadDestFolder 	= "path/"
	//url补丁文件
	urlFilePath        	= "urlList.txt"
	//开启goroutine数量
	NumGroup			= 3
)
var w sync.WaitGroup
func Run() {
	fi, err := os.Open(urlFilePath)
	if err != nil {
		fmt.Printf("Error:%s\n", err)
		return
	}
	defer fi.Close()
	if !IsDir(downloadDestFolder) {
		os.Mkdir(downloadDestFolder,777)
	}
	br := bufio.NewReader(fi)
	c := make(chan string,100000)
	//开启线程
	for i:=0;i<NumGroup;i++{
		go	download(&c)
	}
	for {	
		line, _, err := br.ReadLine()
		if err != nil {
			fmt.Println("readurlcomplete")
			break
		}
		if len(line) > 0 {
			c <- string(line)
		}
	}
	close(c)
	w.Wait()
}

func download(d * chan string){
	w.Add(1)
	defer func(){
		if errs := recover();errs != nil {
			w.Done()
			fmt.Println("err=",errs)
		}
	}()
	var filenames string
	for v := range *d {
		index := strings.LastIndex(v,"/") + 1
		name := v[index:]
		filenames = downloadDestFolder+"\\"+name
		if !isExist(filenames){
			fmt.Println("文件存在:", filenames)
			continue
		}
		url := v
		fmt.Println("下载开始:", url)
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("get err:", err)
			break
		}
		defer res.Body.Close()
		fw,err := os.Create(filenames)
		defer fw.Close()
		if err != nil {
			fmt.Println("create file err:", err)
			break
		}
		r := bufio.NewReader(res.Body)
		if _,err = io.Copy(fw, r); err != nil {
			fmt.Println("write err:", err)
		} else {
			fmt.Println("下载结束:", filenames)
		}
	}
	w.Done()
}
func isExist(path string) bool{
	_,err := os.Stat(path)
	if err != nil {
		boor := os.IsNotExist(err)
		if boor {
			return true
		}
	}
	return false
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}