package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"net/http"	
	"regexp"	
	"time"
	"strconv"	
	"strings"
	"math/rand"
	"github.com/sevenlv007/gostudy/spider/tools"
)

//Spider 入门 正则
//测试成功spider到网页内容
//正则手机号\邮箱\超级链接\身份证\图片地址 抽取成功
//正则里，（）来分切片 

var (
	//手机号	
	//rePhone = `1[3456789]\d{9}`	//手机号的正则
	rePhone = `(1[3456789]\d)(\d{4})(\d{4})`	//改进
	
	//邮箱
	//reEmail = `[1-9]\d{4,}@qq.com`	//非零开始，4位以上，qq.com结尾
	//reEmail = `\w+@\w+.com`	//字符开始，@字符.com结尾
	reEmail = `[\w\.]+@\w+\.[a-z]{2,3}(\.[a-z]{2,3})?`	//改进--0到多个带点的字符开始，@字符.2-3位的小写字母,再.2-3位的小写字母0次或1次
	
	//超链接
	//<a>开头，</a>结尾，其中href后的内容为链接
	reLink = `<a[\s\S]+?href="(http[\s\S]+?)"`

	//身份证
	//第一位1-6，后五位任意数字，然后19或20开通的年份,19后面任意数字，20后面0加任意数字或1加（0-8）
	//月份只能01开头,0开通第二位任意数字，1开头后面只能接012
	//日期为0，12，3开头，0开头后面跟1-9，12后面跟任意数字，3开头只能跟01
	//最后四位，前三位为任意数字，最后一个可以是任意数字或X
	//reId = `[1-6]\d{5}-((19\d{2})|(20((0\d)|(1[0-8]))))-((0[1-9])|(1[012]))-((0[1-9])|([12]\d)|(3[01]))-\d{3}[\dX]`
	reId = `[1-6]\d{5}((19\d{2})|(20((0\d)|(1[0-8]))))((0[1-9])|(1[012]))((0[1-9])|([12]\d)|(3[01]))\d{3}[\dX]`

	//图片
	//<img 开头，src=""开始是图片地址，  > 结尾-可以忽略后面部分
	//reImg = `<img[\s\S]+?src="(http[\s\S]+?)"[\s\S]*?>`
	reImg = `<img.+?src="(http.+?)".*?>`	//下载阻塞-原因不详 少了个”
	//reImg =`<img[\s\S]+?src="(http[\s\S]+?)"` //可以下载
 	//图片包含更多alt信息
	//reImgName = `<img[\s\S]+?alt="([\s\S]+?)"`
	reImgAlt = `<img.+?alt="(.+?)"`

	//reImgWithAlt = `<img.+?src="(.+?)"[\s\S]*?/?>`

	//图片img标签中alt属性
	reAlt = `alt="(.+?)"`
	 
	//图片链接中的名字
	reImgName = `/(\w+\.((jpg)|(jpeg)|(png)|(gif)|(bmp)|(webp)|(swf)|(ico)))`


	//存储img的目录
	imgDir = `D:\go\src\github.com\sevenlv007\gostudy\Spider\img\`
)

//HandleErr 处理错误--
func HandleErr(err error, when string){
	if err != nil {
		fmt.Println(when, err)
		os.Exit(1)			
	}
}

//GetHtml 获取需要spider的网址的GetHtml函数
func GetHtml(url string) string {
	//读入网址	
	//resp, err := http.Get("https://www.haomagujia.com/")	//输入需要spider的网址-电话号码
	resp, _ := http.Get(url)	//输入需要的网址
	//HandleErr(err, `http.Get(url)`)
	defer resp.Body.Close()		//关闭读取，不关闭可能导致网卡断线

	//HandleErr(err, `http.Get`)
	bytes, _ := ioutil.ReadAll(resp.Body)
	
	//读入内容存入string
	html := string(bytes)
	//打印网页内容
	//fmt.Println(html)
	return html
}


//GetPageImgUrls 抓取图片地址的函数,获取页面上所有图片链接
func GetPageImgUrls(url string) []string {
	//获取地址	
	html := GetHtml(url)		//图片
	//fmt.Println(html) 	//测试是否获取到地址

	//使用正则表达式爬取内容
	re := regexp.MustCompile(reImg)	//放入需要爬取的形式 
	allString := re.FindAllStringSubmatch(html, -1) //不知道第二个n是干嘛的
	fmt.Println("spider count:", len(allString))
	
	imgUrls := make([]string, 0)
	for _, x := range allString {
		imgUrl := x[1]
		imgUrls = append(imgUrls, imgUrl)
	}
	return imgUrls		
}

//GetPageImgNameUrls 抓取包含alt信息的图片地址的函数,获取页面上所有图片链接
//内容已由gbk转为utf8
func GetPageImgNameUrls(url string) []string {
	//获取地址	
	html := GetHtml(url)		//图片
	//fmt.Println(html) 	//测试是否获取到地址

	//地址内容 由 gbk转utf8
	bytes := tools.ConvertToString(html, "gbk", "utf8")
	html = string(bytes)
	//fmt.Println(html) 

	//使用正则表达式爬取内容
	re := regexp.MustCompile(reImgAlt)	//放入需要爬取的形式 
	allString := re.FindAllStringSubmatch(html, -1) //不知道第二个n是干嘛的
	fmt.Println("spider count:", len(allString))
	
	imgUrls := make([]string, 0)
	for _, x := range allString {
		imgUrl := x[1]
		imgUrls = append(imgUrls, imgUrl)
	}
	return imgUrls		
}


//GetRandInt 获取随机数 生成[start, end)（含头不含尾）
func GetRandInt(start, end int) int{
	randomMT.Lock()		//上锁
	<-time.After(1 * time.Nanosecond)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ret := start + r.Intn(end-start)
	randomMT.Unlock()	//解锁
	return ret
}

//GetRandName 获取随机文件名_时间戳+随机数
func GetRandName() string {
	timestamp := strconv.Itoa(int(time.Now().UnixNano()))
	randomNum := strconv.Itoa(GetRandInt(100, 1000))
	return timestamp + "_" + randomNum
}


//GetImgNameFromTag 标签中提文件名（含地址）
//使用alt+链接中的名字做文件名
//有alt使用alt作为文件名，没有使用时间戳_随机数做文件名
//参数：
//imgTag图片<img>标签
//dir 目录位置
//suffix 文件名后缀
func GetImgNameFromTag(imgTag, imgUrl, imgDir string) string {

	var filename string
	//获得图片格式
	imgName := GetImgNameFromImgurl(imgUrl)
	suffix := ".jpg"
	if imgName != "" {
		suffix = imgName[strings.LastIndex(imgName, "."):]
	}

	//从imgTag中提取alt
	re := regexp.MustCompile(reAlt)
	rets := re.FindAllStringSubmatch(imgTag, 1)
	if len(rets) > 0 && imgName != "" {
		//首选Alt
		alt := rets[0][1]
		alt =strings.Replace(alt, "：", "_", -1)	//把冒号提成成_,需要区分中文冒号和英文冒号的区别
		filename = alt + imgName
	} else  if imgName != "" {
		//次选文件名
		filename = imgName
	}else{
		//最后选时间戳+随机数
		filename = GetRandName() + suffix
	}
	filename = imgDir + filename 
	return filename
}

//GetImgNameFromImgurl 从imgUrl中抽取图片名称
func GetImgNameFromImgurl(imgUrl string) string {
	re := regexp.MustCompile(reImgName)
	rets := re.FindAllStringSubmatch(imgUrl, -1)
	if len(rets) > 0 {
		return rets[0][1]
	}else {
		return ""
	}
}

//GetPageImgInfos 获取页面上的全部图片信息，链接加文件名
func GetPageImgInfos(url string) []map[string]string {
	html := GetHtml(url)

	html = string(tools.ConvertToString(html, "gbk", "utf8"))
	re := regexp.MustCompile(reImg)
	rets := re.FindAllStringSubmatch(html, -1)
	fmt.Println("Spider countint:", len(rets))

	imgInfos := make([]map[string]string, 0)
	for _, ret := range rets {
		imgInfo := make(map[string]string)
		imgUrl := ret[1]
		imgInfo["url"] = imgUrl
		imgInfo["filename"] = GetImgNameFromTag(ret[0],imgUrl, imgDir)
		imgInfos = append(imgInfos, imgInfo)
	}
	return imgInfos
}


//以reImgalt为目标，爬取html内容
func spiderImg(){
	html := GetHtml("https://www.163.com/")

	//转码
	bytes := tools.ConvertToString(html, "gbk", "utf8")
	html = string(bytes)
	//fmt.Println(html)

	re := regexp.MustCompile(reImgAlt)
	rets := re.FindAllStringSubmatch(html, -1)
	fmt.Println("spider count:", len(rets))
	for _, ret := range rets {
		fmt.Println(ret[1])
	}

}

//以reImg为目标，爬取网页内容，且以Tag内容为文件名
func spiderImgWithAlt() {
	html := GetHtml("https://www.163.com/")

	//转码
	bytes := tools.ConvertToString(html, "gbk", "utf8")
	html = string(bytes)
	//fmt.Println(html)

	re := regexp.MustCompile(reImg)
	rets := re.FindAllStringSubmatch(html, -1)
	fmt.Println("spider count:", len(rets))
	for i, ret := range rets {
		imgTag := ret[0]
		fmt.Println(i, imgTag, "\n", GetImgNameFromTag(imgTag,ret[1],imgDir))
	}
}

//以reImgName为目标，爬取网页内容，且爬取imgurl
func spiderImgName(){
	html := GetHtml("https://www.163.com/")

	//转码
	bytes := tools.ConvertToString(html, "gbk", "utf8")
	html = string(bytes)
	//fmt.Println(html)

	re := regexp.MustCompile(reImgName)
	rets := re.FindAllStringSubmatch(html, -1)
	fmt.Println("spider count:", len(rets))
	for _, ret := range rets {
		imgUrl := ret[1]
		fmt.Println(imgUrl,GetImgNameFromImgurl(imgUrl))

	}
}

//抽取未完成
func dlPage(url string) { 
		//http://www.mmonly.cc/tag/rbsn/index.html
		baseUrl := "http://www.mmonly.cc/tag/rbsn/"
		for i := 1; i < 22; i++ {
			var url string
			if i != 1{
				url = baseUrl + strconv.Itoa(i)+ ".html" //如果有_则需要增加
				fmt.Println("page:", i)
			}else {
				url = baseUrl + "index.html"
				fmt.Println("page:", i)
			}
	
			DownladPageImgs(url)
		}
		wg.Wait()
}



func main() {
	/*
	imgurls := GetPageImgNameUrls("https://www.163.com/")
	for _, imgUrl := range imgurls {
		fmt.Println(imgUrl)
		DownloadImgAsync1(imgUrl)
		//DownloadImg2(imgUrl,GetImgNameFromImgurl(imgUrl))
		//DownloadImgAsync1(imgUrl,GetImgNameFromImgurl(imgUrl))
		//fmt.Println(i,imgurl,"\n", GetImgNameFromTag(imgurl) ,"\n")
	
	}
	
	*/
	/*
	//以下代码已经成功，取名字为alt名
	imgInfos := GetPageImgInfos("http://www.mmonly.cc/tag/rbsn/3.html")
	for _, imgInfoMap := range imgInfos {
		//imgInfos
		//fmt.Println(imgurl)
		
		DownloadImgAsync2(imgInfoMap["url"], imgInfoMap["filename"])
		//DownloadImgAsync1(imgInfoMap["url"])
		//fmt.Println(i,imgInfoMap,"\n", GetImgNameFromTag(imgInfoMap["filename"], imgDir,".jpg") ,"\n")
	}
	*/


	
	/*
	//有bug
	imgurls :=  GetPageImgNameUrls("https://www.163.com/")
	for i, imgUrl := range imgurls {
		fmt.Println(imgUrl)
		//DownloadImg1(imgUrl)
		//DownloadImgAsync1(imgUrl)
		fmt.Println(i,imgUrl,"\n", GetImgNameFromTag(imgUrl, imgDir, ".Jpg") ,"\n")
	}
	*/

	/*spiderImgWithAlt() 
	//以下原版downland的可用，由随即名组成
		imgurls := GetPageImgUrls("https://www.163.com/")
		for _, imgUrl := range imgurls{
		//fmt.Println(imgUrl)
		//DownloadImg(imgUrl)
		DownloadImgAsync1(imgUrl)

	}*/
	
		//http://www.mmonly.cc/tag/rbsn/index.html
	baseUrl := "http://www.mmonly.cc/tag/rbsn/"
	for i := 1; i < 22; i++ {
		var url string
		if i != 1{
			url = baseUrl + strconv.Itoa(i)+ ".html" //如果有_则需要增加
			fmt.Println("page:", i)
		}else {
			url = baseUrl + "index.html"
			fmt.Println("page:", i)
		}

		DownladPageImgs(url)
	}
	wg.Wait()
	
}
