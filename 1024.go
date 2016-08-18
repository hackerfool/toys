package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/cheggaaa/pb.v1"
)

func main() {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		id := fmt.Sprintf("%d", rnd.Intn(270))
		s := []string{"http://www.5caob.com/vod-play-id-", id, "-src-1-num-1.html"}
		URL := strings.Join(s, "")
		resp, err := http.Get(URL)
		if err != nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		str := string(body)
		macName := getMacName(str)
		macURL := getMacURL(str)
		fmt.Println("Download ID:", id, " ", macName, ":", macURL)
		saveMP4(macName, macURL)
		// fmt.Printf("%-56s:%-s\r\n", macName, macURL)
		time.Sleep(time.Duration((rand.Int()%5 + 5)) * time.Second)
	}
}

func getMacName(body string) string {
	return strings.Split(strings.Split(body, "mac_name='")[1], "',")[0]
}

func getMacURL(body string) string {
	return strings.Split(strings.Split(body, "mac_url='")[1], "';")[0]
}

func saveMP4(name, url string) {
	path := []string{"D:\\TEST\\", name, ".mp4"}
	file, err := os.Create(strings.Join(path, ""))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fileSize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	bar := pb.New(fileSize).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.Start()
	bar.ShowSpeed = true
	bar.ShowFinalTime = true
	bar.SetMaxWidth(80)

	writer := io.MultiWriter(file, bar)

	io.Copy(writer, resp.Body)
	bar.Finish()

}
