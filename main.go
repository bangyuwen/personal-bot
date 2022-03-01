package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"github.com/robfig/cron/v3"
)

var lastMsgFG = ""
var lastMsgAAII = ""

func sendMessage(msg string) {
	reqUrl := "https://api.telegram.org/bot1616062145:AAEw9aiOA5Jgo2bjJv7C6iTnUkD3Uu0pMQs/sendMessage?chat_id=-1001500464600&text="
	resp, err := http.Get(reqUrl + url.QueryEscape(msg))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}

func parseAAII(e *colly.HTMLElement) string {
	msg := []string{"AAII"}
	e.ForEach("li", func(_ int, el *colly.HTMLElement) {
		name := el.ChildText(".stat-name")
		date := el.ChildText(".date-label")
		fmt.Println(date)
		val := el.ChildText(".stat-val .val")
		msg = append(msg, name+" "+date+" "+val+"%")
	})
	return strings.Join(msg, "\n")
}

func crawlAAII() {
	c := colly.NewCollector()
	responseMsg := ""

	c.OnHTML(".sidebar-box-block.chart-stat-lastrows", func(e *colly.HTMLElement) {
		responseMsg = parseAAII(e)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.macromicro.me/charts/20828/us-aaii-sentimentsurvey")
	if responseMsg != lastMsgAAII {
		sendMessage(responseMsg)
		lastMsgAAII = responseMsg
	}
}

func parseFG(e *colly.HTMLElement) string {
	msg := []string{"CNN Fear & Greed Index"}
	updated_at := e.ChildText("#needleAsOfDate")
	e.ForEach("div#needleChart li", func(_ int, el *colly.HTMLElement) {
		msg = append(msg, strings.ReplaceAll(el.Text, "Fear & Greed ", ""))
	})
	msg = append(msg, updated_at)
	return strings.Join(msg, "\n")
}

func crawlFG() {
	c := colly.NewCollector()
	responseMsg := ""

	c.OnHTML("div.modContent.feargreed", func(e *colly.HTMLElement) {
		responseMsg = parseFG(e)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://money.cnn.com/data/fear-and-greed/")
	if responseMsg != lastMsgFG {
		sendMessage(responseMsg)
		lastMsgFG = responseMsg
	}
}

func main() {
	c := cron.New()
	c.AddFunc("@hourly", crawlFG)
	c.AddFunc("@hourly", crawlAAII)
	fmt.Println("Start Cron Crawl")
	fmt.Println(c.Entries())
	c.Start()
	select {}
}
