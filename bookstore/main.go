package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/HaTuanZ/crawling/common"
	"github.com/gocolly/colly"
)

type Item struct {
	Title string `json:"title"`
	Price string `json:"price"`
}

var url = "https://nhanam.vn/van-hoc-hien-dai?q=collections:3204401&page=1&view=grid"

func main() {
	var items []Item
	c := colly.NewCollector()

	c.OnHTML("body > main > div > div > div.bg_collection.section > div > div > div.category-products.products > section > div", func(h *colly.HTMLElement) {
		h.ForEach(".product-col", func(i int, h *colly.HTMLElement) {
			items = append(items, Item{
				Title: h.ChildText(" div > form > div.info-product > h3 > a"),
				Price: h.ChildText("div > form > div.info-product > div"),
			})
		})
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.Visit(url)

	content, err := json.Marshal(items)

	if err != nil {
		fmt.Println(err.Error())
	}

	os.WriteFile(common.GetCurrentPath()+"/bookstore/books.json", content, 0664)

}
