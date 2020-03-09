package main

import (
	"fmt"
	"io/ioutil"
	"net/http"	
	//"time"
	//"strconv"	
	"sync"
		
)
var (
	//并发管道	管道数10个
	chSem = make(chan int, 5)
	//图片下载等待组
	downloadWG sync.WaitGroup
	//互斥锁
	randomMT sync.Mutex
	
)

//DownloadImg1 同步下载图片的函数-随机文件名
func DownloadImg1(url string) {
	resp, _:= http.Get(url)
	//HandleErr(err, `http.Get(url)`)
	defer resp.Body.Close()

	imgbytes, _ := ioutil.ReadAll(resp.Body)
	filename := `D:\go\src\github.com\sevenlv007\gostudy\Spider\img\`+GetRandName()+".jpg"
	err := ioutil.WriteFile(filename, imgbytes, 0644)
	if err == nil {
		fmt.Println("Download Succeed")
	}else{
		fmt.Println("Download failed")
	}
}

//DownloadImg2 同步下载图片的函数-传入文件名
func DownloadImg2(url, filename string) {
	fmt.Println("Downloading Img....")
	resp, _ := http.Get(url)
	//HandleErr(err, `http.Get(url)`)

	defer resp.Body.Close()

	imgBytes, _ := ioutil.ReadAll(resp.Body)
	//filename = `D:\go\src\github.com\sevenlv007\gostudy\Spider\img\`+strconv.Itoa(int(time.Now().UnixNano()))+".jpg"
	
	err := ioutil.WriteFile(filename, imgBytes, 0644)
	if err == nil {
		fmt.Println(filename+" Download Succeed")
	}else{
		fmt.Println(filename+" Download failed")
	}
}



//DownloadImgAsync1 异步下载图片函数-不传入文件名
func DownloadImgAsync1(url string) {
	downloadWG.Add(1)	//不加等待，主线程会死
	go func(){
		chSem <- 123
		DownloadImg1(url)
		<-chSem
		downloadWG.Done()
	}()

	downloadWG.Wait()
}


//DownloadImgAsync2 异步下载图片函数-传入文件名
func DownloadImgAsync2(url, filename string) {
	downloadWG.Add(1)	//不加等待，主线程会死
	go func(){
		chSem <- 123
		DownloadImg2(url, filename)
		<-chSem
		downloadWG.Done()
	}()
	downloadWG.Wait()
}

