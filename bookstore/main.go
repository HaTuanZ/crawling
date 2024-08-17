package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sync"

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
	pageNumber := 13
	chunkNumber := 3
	workerNumber := int(math.Ceil(float64(pageNumber) / float64(chunkNumber)))
	workers := make([]chan string, workerNumber)
	wg := new(sync.WaitGroup)

	// c.OnHTML(".page-item:nth-last-child(2) .page-link", func(h *colly.HTMLElement) {
	// 	re := regexp.MustCompile("[^0-9]")
	// 	var err error

	// 	pageNumber, err = strconv.Atoi(re.ReplaceAllString(h.Text, "*"))
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// })

	for i := 0; i < workerNumber; i++ {
		chunk := chunkNumber * (i + 1)
		workers[i] = make(chan string, chunkNumber)
		for k := i*chunkNumber + 1; k <= chunk && k <= pageNumber; k++ {
			workers[i] <- fmt.Sprintf("https://nhanam.vn/van-hoc-hien-dai?q=collections:3204401&page=%d&view=grid", k)
		}
		close(workers[i])
	}

	for i := 0; i < workerNumber; i++ {
		wg.Add(1)
		go func(i int) {
			c := colly.NewCollector(colly.Async(true))

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

			for v := range workers[i] {
				c.Visit(v)
			}
			c.Wait()
			wg.Done()
		}(i)
	}

	wg.Wait()

	content, err := json.Marshal(items)

	if err != nil {
		fmt.Println(err.Error())
	}

	os.WriteFile(common.GetCurrentPath()+"/bookstore/books.json", content, 0664)

}
