package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	now := time.Now()

	fmt.Printf("Start: %s\n", now)

	ts := now.Format("20060102150405")

	url := "https://netmall.hardoff.co.jp/cate/000100050001/"
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.79"

	client := &http.Client{Timeout: 30 * time.Second}

	file, err := os.Create(ts + ".csv")
	if err != nil {
		log.Fatalf("faild: %s", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	err = w.Write([]string{"SOLDOUT", "BRAND", "TITLE", "TAG", "PRICE", "URL"})
	if err != nil {
		log.Fatalf("faild: %s", err)
	}

	for {
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Add("User-Agent", userAgent)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("faild request: %s", err)
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Printf("faild body: %s", err)
		}

		doc.Find(".itemcolmn_item").Each(func(i int, s *goquery.Selection) {
			soldout := s.Find(".soldout-text").First().Text()
			brand := s.Find(".item-brand-name").First().Text()
			title := s.Find(".item-name").First().Text()
			tag := s.Find(".item-code").First().Text()
			price := s.Find(".item-price-en").First().Text()
			price = strings.ReplaceAll(price, "\t", "")
			price = strings.ReplaceAll(price, "\n", "")
			price = strings.ReplaceAll(price, ",", "")
			price = strings.ReplaceAll(price, "円", "")
			a := s.Find("a")
			val, exists := a.Attr("href")
			if exists {
				err := w.Write([]string{soldout, brand, title, tag, price, val})
				if err != nil {
					log.Fatalf("faild: %s", err)
				}
			}
		})

		next := doc.Find(".next").Length()

		if next > 0 {
			doc.Find(".next").Each(func(i int, s *goquery.Selection) {
				val, exists := s.Attr("href")
				fmt.Printf("Next page: %s\n", val)
				if exists {
					url = val
				}
			})
		} else {
			break
		}

		resp.Body.Close()

		time.Sleep(time.Second * 2)
	}

	fmt.Printf("End: %s\n", time.Now())
}
