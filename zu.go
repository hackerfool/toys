package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const KEY = "e7cb366445285ea602afe57d33678436"
const WEBKEY = "1914013c93d688a460b6dc37371afafc"

//{"status":"1","info":"OK","infocode":"10000","count":"1","geocodes":[{"formatted_address":"上海市浦东新区世纪大道","province":"上海市","citycode":"021","city":"上海市","district":"浦东新区","township":[],"neighborhood":{"name":[],"type":[]},"building":{"name":[],"type":[]},"adcode":"310115","street":[],"number":[],"location":"121.524565,31.229925","level":"道路"}]}
type geoResp struct {
	Status   string
	Info     string
	InfoCode string
	Count    string
	GeoCodes []geoCodes
}

type geoCodes struct {
	Formatted_address string
	Province          string
	CityCode          string
	City              string
	Disctrict         string
	Location          string
	Level             string
}

type directionResp struct {
	Status   string
	Info     string
	InfoCode string
	Count    string
	Route    route
}

type route struct {
	Origin      string
	Destination string
	Distance    string
	Taxi_cost   string
	Transits    []transit
}

type transit struct {
	Cost     string
	Duration string
}

type apartment struct {
	Name      string
	Address   string
	Longitude string
	Latitude  string
}

var mapHtml = `<!doctype html>
			<html>
			<head>
				<meta charset="utf-8">
				<meta http-equiv="X-UA-Compatible" content="IE=edge">
				<meta name="viewport" content="initial-scale=1.0, user-scalable=no, width=device-width">
				<title>公交路线规划－使用默认样式</title>
				<link rel="stylesheet" href="http://cache.amap.com/lbs/static/main.css?v=1.0"/>
				<script type="text/javascript"
						src="http://webapi.amap.com/maps?v=1.3&key=%s"></script>
				<style type="text/css">
					#panel {
						position: absolute;
						background-color: white;
						max-height: 80%%;
						overflow-y: auto;
						top: 10px;
						right: 10px;
						width: 250px;
						border: solid 1px silver;
					}
				</style>
			</head>
			<body>
			<div id="mapContainer"></div>
			<div id="panel">
			</div>
			<script type="text/javascript">
				var map = new AMap.Map("mapContainer", {
					resizeEnable: true,
					center: [121.443133, 31.280432],//地图中心点
					zoom: 13 //地图显示的缩放级别
				});
				/*
				* 调用公交换乘服务
				*/
				//加载公交换乘插件
				AMap.service(["AMap.Transfer"], function() {
					var transOptions = {
						map: map,
						city: '%s',
						panel:'panel',                            //公交城市
						//cityd:'乌鲁木齐',
						policy: AMap.TransferPolicy.LEAST_TIME //乘车策略
					};
					//构造公交换乘类
					var trans = new AMap.Transfer(transOptions);
					//根据起、终点坐标查询公交换乘路线
					trans.search([%s],[%s], function(status, result){
					});
				});
			</script>
			<script type="text/javascript" src="http://webapi.amap.com/demos/js/liteToolbar.js"></script>
			</body>
			</html>`

func DistanceHandle(w http.ResponseWriter, req *http.Request) {
	html := `<html><body>%s</body></html>`
	body := ""
	req.ParseForm()
	destAddress, destLocation := geo(req.Form["dest"][0])
	for _, value := range req.Form["org"] {
		address, location := geo(value)
		ways := distance(location, destLocation, "上海市")
		body += fmt.Sprintf("%s(%s) -- %s(%s)</br>%s<a href=\"./map?org=%s&dest=%s&city=%s\">公交地图</a></br>", address, location, destAddress, destLocation,
			ways, location, destLocation, "上海市")
	}
	html = fmt.Sprintf(html, body)
	io.WriteString(w, html)
}

func MapHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	html := fmt.Sprintf(mapHtml, WEBKEY, req.Form["city"][0], req.Form["org"][0], req.Form["dest"][0])
	io.WriteString(w, html)
}

func main() {
	// http.HandleFunc("/distance", DistanceHandle)
	// http.HandleFunc("/map", MapHandle)
	// log.Fatal(http.ListenAndServe(":12345", nil))
	getApartment("http://sh.lianjia.com/xiaoqu/5011000017674.html")
}

func geo(address string) (string, string) {
	url := fmt.Sprintf("http://restapi.amap.com/v3/geocode/geo?address=%s&output=json&key=%s", address, KEY)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	// body := fmt.Sprintf("%s", robots)
	js := &geoResp{}
	err = json.Unmarshal(robots, js)
	if err != nil {
		log.Fatal(err)
	}
	//fix 多个结果
	return js.GeoCodes[0].Formatted_address, js.GeoCodes[0].Location
}

func distance(org string, dest string, city string) string {
	//http://restapi.amap.com/v3/direction/transit/integrated?origin=116.481499,39.990475&destination=116.465063,39.999538&city=010&output=xml&key=<用户的key>
	url := fmt.Sprintf("http://restapi.amap.com/v3/direction/transit/integrated?origin=%s&destination=%s&city=%s&output=json&key=%s", org, dest, city, KEY)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	// body := fmt.Sprintf("%s", robots)
	js := &directionResp{}
	err = json.Unmarshal(robots, js)
	if err != nil {
		log.Fatal(err)
	}
	body := ""
	for i, x := range js.Route.Transits {
		sec, _ := strconv.Atoi(x.Duration)
		body += fmt.Sprintf("方案%d:需要%d分</br>", i+1, sec/60)
	}
	return body
}

func getApartment(url string) {
	document, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	attr, ok := document.Find("#actshowMap_xiaoqu").Attr("xiaoqu")
	address,_ := document.Find("body div.wrapper div.nav-container.detail-container section div.res-top.clear div.title.fl span span.adr").Attr("title")
	if !ok {
		return
	}
	result := &apartment{}
	attrs := strings.Split(attr, ",")
	result.Longitude = attrs[0][1:]
	result.Latitude = attrs[1]
	result.Name = attrs[2][2:len(attrs[2])-2]
	result.Address = address
	fmt.Println(result)
}