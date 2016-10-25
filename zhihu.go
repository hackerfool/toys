package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
)

//test for rebase
//test for rebase 2

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
	request.Header.Set("Cookie", `q_c1=746738a68faf490e869c8edf240cc728|1473842100000|1473842100000; _xsrf=4b80377c20555b6b2f9fbdfdc88bbf06; l_cap_id="YmZjY2MzZmFmY2QyNDc5ZDgxYzYxZGQzYTY3NzBiOWI=|1473842100|916a4c17d9287b116e938976c7a4fadff3f2c653"; cap_id="ZjQ5ODU3MDkxMzhhNDkyN2IzNTlkNGNmZmI3N2FlZGU=|1473842100|4801fdaf1c3e8900b57cc4ff10adff2c69d08104"; _za=8faf509f-a4af-4e29-9c7f-8d05ea05c6da; d_c0="ADDAZhoJigqPToHwAUdAJ7vTgS8ABdWGI8Q=|1473842101"; _zap=167ff54b-3dba-412d-a4a3-a9e214ab7968; n_c=1; s-q=%E6%94%AF%E4%BB%98%E5%AE%9D%20%E5%88%B0%E4%BD%8D; s-i=1; sid=tcajbglg; __utmt=1; a_t="2.0AACAG6EdAAAXAAAAKaAAWAAAgBuhHQAAADDAZhoJigoXAAAAYQJVTcKYAFgAiEFGGQe5AovD6OerspDKaHRDIO26fqHdkGH6DTQEmHwe05CHyxtutg=="; z_c0=Mi4wQUFDQUc2RWRBQUFBTU1CbUdnbUtDaGNBQUFCaEFsVk53cGdBV0FDSVFVWVpCN2tDaThQbzU2dXlrTXBvZEVNZzdR|1473844009|dea2565924d26db92162db21b001638ee1b675c8; __utma=51854390.604875030.1473842155.1473842155.1473842155.1; __utmb=51854390.6.10.1473842155; __utmc=51854390; __utmz=51854390.1473842155.1.1.utmcsr=zhihu.com|utmccn=(referral)|utmcmd=referral|utmcct=/search; __utmv=51854390.100-1|2=registration_date=20130829=1^3=entry_date=20130829=1; a=; _za=16974060-84bb-4fa0-ba14-fdb86fe5e886; d_c0="ADBA0CxtoQmPTmo5cFtch-v75PS4jUZ5PeM=|1458231969"; _zap=5808caf8-b490-4299-ae8f-6dcdd1055c8b; q_c1=1a4596f2ae4d4ee2af35ad24d2f717f1|1471853877000|1471853877000; l_cap_id="MTcwZDViMzMxZGJlNGJmZmFmOTA0YjZiMmJhYjM1YzI=|1473326457|600254870fe296387993505206c28e6f4ace656a"; cap_id="ZDc5MWM4ZGNkZmQ4NGY5MWI4MDVmNDJiMmQzMzJlMDc=|1473326457|cd4c7515c1a155732a1a13e109241801fdc245da"; _xsrf=d9eda4b69ca818fa93cba72c0c89efd4; __utmt=1; __utma=51854390.1821037221.1465872812.1473754113.1473844880.27; __utmb=51854390.10.10.1473844880; __utmc=51854390; __utmz=51854390.1473754113.26.19.utmcsr=t.co|utmccn=(referral)|utmcmd=referral|utmcct=/YO2xECvADn; __utmv=51854390.100-1|2=registration_date=20130829=1^3=entry_date=20130829=1; a_t="2.0AACAG6EdAAAXAAAA4aMAWAAAgBuhHQAAADBA0CxtoQkXAAAAYQJVTYO6-FcAT1l-wRNOXGhNOX93by5S6ROrsFsl_r2FuqKqDOtqQLTC10EWtv-mUQ=="; z_c0=Mi4wQUFDQUc2RWRBQUFBTUVEUUxHMmhDUmNBQUFCaEFsVk5nN3I0VndCUFdYN0JFMDVjYUUwNWYzZHZMbExwRTZ1d1d3|1473844961|890393590596cf9220c446925d964ecb90f98e1a`)
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

	xsrf, _ := followBody.Find(`input[name="_xsrf"]`).Attr("value")

	date, _ := followBody.Find("div.zh-general-list.clearfix").Attr("data-init")
	js, err := simplejson.NewJson([]byte(date))
	//{"params": {"offset": 0, "order_by": "created", "hash_id": "c8a506f7961e94bb0ef3f243b198be70"}, "nodename": "ProfileFolloweesListV2"}
	hashID := js.Get("params").Get("hash_id").MustString()
	// fmt.Println(hashID)
	for i := 0; i < person.followeesNum; i += 20 {
		str := strings.Join([]string{"method=next&params=%7B%22offset%22%3A", strconv.Itoa(i), "%2C%22order_by%22%3A%22created%22%2C%22hash_id%22%3A%22", hashID, "%22%7D"}, "")
		req, err := http.NewRequest("POST", "https://www.zhihu.com/node/ProfileFolloweesListV2", strings.NewReader(str))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Referer", URL)
		req.Header.Set("X-Xsrftoken", xsrf)

		if err != nil {
			log.Fatal(err)
		}
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		web, err := ioutil.ReadAll(res.Body)
		if err != nil {
			// handle error
		}

		js, err = simplejson.NewJson(web)
		msg := js.Get("msg").MustStringArray()
		msgstr := strings.NewReader(strings.Join(msg, ""))
		doc, err := goquery.NewDocumentFromReader(msgstr)
		doc.Find("div.zm-profile-card.zm-profile-section-item.zg-clear.no-hovercard").Each(func(i int, dom *goquery.Selection) {
			s := dom.Find("a.zm-item-link-avatar")
			name, _ := s.Attr("title")
			url, _ := s.Attr("href")
			fmt.Println(name, "--", url)
		})
	}
}
