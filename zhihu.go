package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
)

//Person 知乎用户信息
type Person struct {
	name         string
	location     string
	business     string
	sex          int
	url          string
	avatar       string
	getAgree     int
	getThanks    int
	ask          int
	answer       int
	posts        int
	collections  int
	logs         int
	followeesNum int
	followersNum int
	followees    map[string]Person
	followers    map[string]Person
}

var client *http.Client

func main() {
	url := "https://www.zhihu.com/people/zhumo0.0"
	//url := "https://www.zhihu.com/people/wan-nian"
	cookieJar, _ := cookiejar.New(nil)
	client = &http.Client{Jar: cookieJar}
	readPeople(url)
}

func readPeople(url string) {

	request, err := http.NewRequest("GET", url, strings.NewReader(""))
	request.Header.Set("Cookie", `d_c0="AJCAFHjAiAqPThKVF-mjCBfcWA8mMqvxBm8=|1473755951"; _za=38c06cb5-7f0c-445c-a67a-c4d59b6493cf; _xsrf=50f89e23ebe1b959838c01208a6ede8f; _zap=7641b9c6-8f37-4179-b70a-2d59f2b3e8c3; q_c1=3c9b1116b67d4779805d24ee4c6ed9c2|1473773979000|1473773979000; l_cap_id="NjIyMGZlNDU0MWVmNDVhZWExYjc3ZTMwMWNmYmRmYWI=|1473774826|0f2ed9e638b8e47c82fc68d07019eb8e6131c1ec"; cap_id="Zjk5OTQ4NWZmMGZkNDM1YWI3YWY0ZDYwMmI5NzIzOTA=|1473774826|354fe08c4ca32589c6f889bcb66f3ff51d5c4200"; n_c=1; a_t="2.0AACAG6EdAAAXAAAAULD_VwAAgBuhHQAAAJCAFHjAiAoXAAAAYQJVTe6R_1cAorqoivJAOegGeKpl1glWAhOM9JgertH_e2ksI7FXyEf0qsnb50PTOg=="; z_c0=Mi4wQUFDQUc2RWRBQUFBa0lBVWVNQ0lDaGNBQUFCaEFsVk43cEhfVndDaXVxaUs4a0E1NkFaNHFtWFdDVllDRTR6MG1B|1473782608|5e2da97cd4fc22e4a8e46801df0da73d7b78ba0b; __utmt=1; __utma=51854390.930195786.1473773031.1473773031.1473782610.2; __utmb=51854390.3.9.1473782612648; __utmc=51854390; __utmz=51854390.1473782610.2.2.utmcsr=zhihu.com|utmccn=(referral)|utmcmd=referral|utmcct=/people/zhumo0.0; __utmv=51854390.100-1|2=registration_date=20130829=1^3=entry_date=20130829=1`)
	// body, err := goquery.NewDocument(url)
	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	body, err := goquery.NewDocumentFromResponse(response)
	// fmt.Println(body.Text())
	if err != nil {
		log.Fatal(err)
	}
	person := &Person{url: url}
	getName(body, person)
	getNav(body, person)
	getFollow(body, person)
	// fmt.Println(person)
}

func getNav(body *goquery.Document, person *Person) {
	navAttr := map[int]string{}
	body.Find(".profile-navbar").Find("a.item").Each(func(i int, d *goquery.Selection) {
		navAttr[i] = d.Find("span.num").Text()
	})
	person.ask, _ = strconv.Atoi(navAttr[1])
	person.answer, _ = strconv.Atoi(navAttr[2])
	person.posts, _ = strconv.Atoi(navAttr[3])
	person.collections, _ = strconv.Atoi(navAttr[4])
	person.logs, _ = strconv.Atoi(navAttr[5])
}

func getName(body *goquery.Document, person *Person) {
	person.name = body.Find("span.name").Eq(0).Text()
	person.location, _ = body.Find("span.location").Attr("title")
	person.business, _ = body.Find("span.business.item").Attr("title")
	if body.Has("i.icon.icon-profile-male").Size() != 0 {
		person.sex = 0
	} else {
		person.sex = 1
	}

	// Avatar Avatar--l 处理用户头像
	avatarURL, exits := body.Find("img.Avatar.Avatar--l").Attr("src")
	if exits == true {
		person.avatar = strings.Join(strings.Split(avatarURL, "_l"), "")
	}
}

func getFollow(body *goquery.Document, person *Person) {
	follow := body.Find("div.zm-profile-side-following").Find("a.item")
	followeesHref, _ := follow.Eq(0).Attr("href")
	followessNum := follow.Eq(0).Find("strong").Text()
	// followersHref, _ := follow.Eq(1).Attr("href")
	// followersNum := follow.Eq(1).Find("strong").Text()
	person.followeesNum, _ = strconv.Atoi(followessNum)
	// person.followersNum, _ = strconv.Atoi(followersNum)
	getUserFollowees(followeesHref, person)
}

func getUserFollowees(URL string, person *Person) {
	//获取关注了的用户
	Href := strings.Join([]string{"https://www.zhihu.com", URL}, "")
	request, err := http.NewRequest("GET", Href, nil)
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	followBody, _ := goquery.NewDocumentFromResponse(response)
	// followBody.Find("div#zh-profile-follows-list").Find("div.zm-profile-card.zm-profile-section-item.zg-clear.no-hovercard").Each(func(i int, dom *goquery.Selection) {
	// 	s := dom.Find("a.zm-item-link-avatar")
	// 	name, _ := s.Attr("title")
	// 	url, _ := s.Attr("href")
	// 	fmt.Println(name, "--", url)
	// })
	date, _ := followBody.Find("div.zh-general-list.clearfix").Attr("data-init")
	js, err := simplejson.NewJson([]byte(date))
	//{"params": {"offset": 0, "order_by": "created", "hash_id": "c8a506f7961e94bb0ef3f243b198be70"}, "nodename": "ProfileFolloweesListV2"}
	hashID := js.Get("params").Get("hash_id").MustString()
	// fmt.Println(hashID)
	for i := 0; i < person.followeesNum; i += 20 {
		str := strings.Join([]string{"method=next&params=%7B%22offset%22%3A", strconv.Itoa(i), "%2C%22order_by%22%3A%22created%22%2C%22hash_id%22%3A%22", hashID, "%22%7D"}, "")
		req, err := http.NewRequest("POST", "https://www.zhihu.com/node/ProfileFolloweesListV2", strings.NewReader(str))
		if err != nil {
			log.Fatal(err)
		}
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		div, err := goquery.NewDocumentFromResponse(res)
		fmt.Println(div.Text())
	}
}
