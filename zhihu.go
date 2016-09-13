package main

import (
	"fmt"

	"github.com/opesun/goquery"
)

func main() {
	URL := "https://www.zhihu.com/people/zhumo0.0"
	x, err := goquery.ParseUrl(URL)
	if err != nil {
		panic(err)
	}
	y := x.Find("div.zm-profile-side-following")
	(y.Find("a.item")).Find("strong").Print()
	// x.Find("span.zg-gray-normal").First().Text()
	fmt.Println("---")
}
