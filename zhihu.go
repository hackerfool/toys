package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	URL := "https://www.zhihu.com/people/zhumo0.0"
	readPeople(URL)
}

func readPeople(url string) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, strings.NewReader(""))
	//request.Header.Set("Referer", "https://www.zhihu.com/")
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}
