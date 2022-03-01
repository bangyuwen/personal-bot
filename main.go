package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"github.com/robfig/cron/v3"
)

func parse(e *colly.HTMLElement) string {
	msg := []string{"CNN Fear & Greed Index"}
	updated_at := e.ChildText("#needleAsOfDate")
	e.ForEach("div#needleChart li", func(_ int, el *colly.HTMLElement) {
		msg = append(msg, strings.ReplaceAll(el.Text, "Fear & Greed ", ""))
	})
	msg = append(msg, updated_at)
	return strings.Join(msg, "\n")
}

func crawl() {
	c := colly.NewCollector()
	responseMsg := ""

	c.OnHTML("div.modContent.feargreed", func(e *colly.HTMLElement) {
		responseMsg = parse(e)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://money.cnn.com/data/fear-and-greed/")
	reqUrl := "https://api.telegram.org/bot1616062145:AAEw9aiOA5Jgo2bjJv7C6iTnUkD3Uu0pMQs/sendMessage?chat_id=-1001500464600&text="
	resp, err := http.Get(reqUrl + url.QueryEscape(responseMsg))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}

func main() {
	// crawl()
	c := cron.New()
	c.AddFunc("@hourly", crawl)
	fmt.Println("Start Cron Crawl")
	fmt.Println(c.Entries())
	c.Start()
	select {}
}
