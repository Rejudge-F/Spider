package main

import (
	"fmt"
	"net/http"

	"os"
	"regexp"
	"strconv"
)

func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()
	buf := make([]byte, 4*1024)
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		result += string(buf[:n])
	}
	return
}

func GetJoy(url string, ch chan int, th int) (result string, err error) {
	result, err1 := HttpGet(url)
	if err1 != nil {
		err = err1
		return
	}
	str := "<h1 class=\"f18\"><a href=\"" + url + "\" title=\""
	str += "(?s:(.*?))\">"

	re := regexp.MustCompile(str)

	if re == nil {
		fmt.Println("complile error")
		return
	}
	//fmt.Println("adfadfadsfasdfd", result)
	titles := re.FindAllStringSubmatch(result, 1)
	var title string
	for _, data := range titles {
		title = data[1]
		break
	}

	//fmt.Println("title = ", title)

	file, err1 := os.Create("./file/" + title + ".txt")
	defer file.Close()
	if err1 != nil {
		err = err1
		return
	}
	str = "<div class=\"con-txt\">" + "(?s:(.*?))</div>"
	//fmt.Println(str)
	re1 := regexp.MustCompile(str)
	if re1 == nil {
		fmt.Println("regexp failed")
		return
	}
	joys := re1.FindAllStringSubmatch(result, -1)
	for _, data := range joys {
		joy := data[1]
		file.WriteString(joy[12:])
	}
	ch <- th
	return
}

func SpiderPage(i int, page chan int) {
	url := "https://m.pengfue.com/index_" + strconv.Itoa(i) + ".html"
	fmt.Printf("正在爬取第%d页，url=%s\n", i, url)
	result, err := HttpGet(url)
	channel := make(chan int)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	//fmt.Println(result)
	re := regexp.MustCompile(`<h1 class="f18"><a href="(?s:(.*?))"`)
	if re == nil {
		fmt.Println("re err = ", re)
		return
	}
	joyUrls := re.FindAllStringSubmatch(result, -1)

	for pos, data := range joyUrls {
		go GetJoy(data[1], channel, pos)
	}

	for p := 0; p < len(joyUrls); p++ {
		fmt.Printf("第%d页的第%d个笑话已经爬取完成...\n", i, <-channel)
	}

	page <- i
}

func DoWork(start, end int) {
	fmt.Printf("准备爬取%d -> %d的页面\n", start, end)
	page := make(chan int)
	for i := start; i <= end; i++ {
		go SpiderPage(i, page)
	}
	for i := start; i <= end; i++ {
		fmt.Printf("第%d页完成！\n", <-page)
	}
}

func main() {

	var start, end int
	fmt.Println("请输入起始页（>=1）：")
	fmt.Scan(&start)
	fmt.Println("请输入终止页（>=起始页)")
	fmt.Scan(&end)

	DoWork(start, end)
}
