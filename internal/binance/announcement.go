package binance

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

/**
币安公告里包含：

1. 上新币
2. 质押挖矿
*/

// HTML curl https://www.binance.com/zh-CN/feed/profile/binance_announcement

const urlStr = "https://www.binance.com/zh-CN/feed/profile/binance_announcement"

func GetList() {
	// 设置代理服务器 生产服务器可去掉
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

	doc.Find("div.feed-content-text").Each(func(i int, s *goquery.Selection) {
		_ = GetPledgeToken(s)
	})

}

func GetNewToken(s *goquery.Selection) error {
	pattern := `币安上市(.*?)并为其添加种子标签`
	re := regexp.MustCompile(pattern)

	var eventList []Event
	a := s.Find("a")
	text := a.Text()
	if href, exists := a.Attr("href"); exists {
		if !strings.Contains(text, "币安新币挖矿上线第") {
			return nil
		}
		match := re.FindStringSubmatch(text)
		if match == nil {
			return nil
		}
		eventList = append(eventList, Event{
			EventTypeName: "新币上市",
			SubTokenName:  match[1],
			MsgUrl:        "https://www.binance.com" + href,
			PushTime:      "",
			Msg:           text,
		})
	}
	fmt.Printf("%+v", eventList)
	return nil
}

// 质押挖矿
func GetPledgeToken(s *goquery.Selection) error {
	pattern := `使用.*?挖矿(.*?)）`
	re := regexp.MustCompile(pattern)

	var eventList []Event
	a := s.Find("a")
	text := a.Text()
	if href, exists := a.Attr("href"); exists {
		if !strings.Contains(text, "币安新币挖矿上线第") {
			return nil
		}
		match := re.FindStringSubmatch(text)
		if match == nil {
			return nil
		}
		eventList = append(eventList, Event{
			EventTypeName: "质押挖币",
			SubTokenName:  match[1],
			MsgUrl:        "https://www.binance.com" + href,
			PushTime:      "",
			Msg:           text,
		})
	}
	fmt.Printf("%+v", eventList)
	return nil
}

type Event struct {
	EventTypeName string
	SubTokenName  string // 质押获得币
	MsgUrl        string // 信息链接
	PushTime      string // 发布时间
	Msg           string // 发布信息
}
