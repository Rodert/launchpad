package binance

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
)

// HTML curl https://www.binance.com/zh-CN/feed/profile/binance_announcement

const urlStr = "https://www.binance.com/zh-CN/feed/profile/binance_announcement"

func GetList() {
	// 设置代理服务器
	proxyURL, err := url.Parse("http://127.0.0.1:7890")
	if err != nil {
		log.Fatal(err)
	}

	// 创建代理客户端
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	// 使用代理客户端发送 HTTP 请求
	res, err := httpClient.Get(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	//doc, err := goquery.NewDocument(urlStr)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Print(err)
	}
	var txt = []string{}

	doc.Find("div.feed-content-text").Each(func(i int, s *goquery.Selection) {
		a := s.Find("a")
		if href, exists := a.Attr("href"); exists {
			fmt.Println("https://www.binance.com" + href)
		}
		txt = append(txt, a.Text())
		fmt.Println(txt)
		fmt.Println("\n\n\n--------\n\n\n")
	})

}
