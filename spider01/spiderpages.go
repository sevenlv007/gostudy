package main

import (
	"sync"


)
var wg sync.WaitGroup

//DownladPageImgs 函数用来爬去多页图片 
func DownladPageImgs(url string){
	//imgUrls := GetHtml("https://www.duotoo.com/zt/rbmn/index.html")
	imgInfos := GetPageImgInfos(url)
	for _, imgInfoMap := range imgInfos {
		//imgInfos
		//fmt.Println(imgurl)
		
		DownloadImgAsync3(imgInfoMap["url"], imgInfoMap["filename"],&wg)
		//DownloadImgAsync1(imgInfoMap["url"])
		//fmt.Println(i,imgInfoMap,"\n", GetImgNameFromTag(imgInfoMap["filename"], imgDir,".jpg") ,"\n")
	}
		


}

//DownloadImgAsync3 用于DownladPageImgs的异步下载方式
func DownloadImgAsync3(url, filename string, wg *sync.WaitGroup) {
	wg.Add(1)	//不加等待，主线程会死
	go func(){
		chSem <- 123
		DownloadImg2(url, filename)
		<-chSem
		wg.Done()
	}()
	
}