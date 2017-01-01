package models

import (
	// "fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/astaxie/beego"
)

type HomeSong struct {
	Url   string
	Title string
	Name  string
}

type HomeList struct {
	Id    string
	Title string
	Name  string
	Views int
}

const (
	NAME = "音乐查询工具"
)

var SongRecord [4]HomeSong
var ListRecord [3]HomeList

func HttpGet(url string) (body []byte) {
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	reqest.Header.Set("Cookie", "appver=1.5.0.75771")
	reqest.Header.Set("Referer", "http://music.163.com/")

	response, err := client.Do(reqest)
	defer response.Body.Close()

	if err != nil {
		return
	}

	if response.StatusCode == 200 {
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return
		}
	}
	return
}

func HttpPost(url, data string) (body []byte) {
	client := &http.Client{}
	reqest, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return
	}
	reqest.Header.Set("Cookie", "appver=1.5.0.75771")
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqest.Header.Set("Referer", "http://music.163.com/")

	response, err := client.Do(reqest)
	defer response.Body.Close()

	if err != nil {
		return
	}

	if response.StatusCode == 200 {
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return
		}
	}
	return
}

func GetId(url string) (id string) {
	reg, _ := regexp.Compile(`[1-9]([0-9]{3,11})`)
	id = reg.FindString(url)
	return
}

func AddSong(song HomeSong) {
	for i := 0; i < 4; i++ {
		if song == SongRecord[i] {
			return
		}
	}
	var old [4]HomeSong
	old[0] = song
	for i := 0; i < 3; i++ {
		old[i+1] = SongRecord[i]
	}
	SongRecord = old
}

func AddList(list HomeList) {
	for i := 0; i < 3; i++ {
		if list == ListRecord[i] {
			return
		}
	}
	var old [3]HomeList
	old[0] = list
	for i := 0; i < 2; i++ {
		old[i+1] = ListRecord[i]
	}
	ListRecord = old
}

func StrGetBetween(str, start, end string) string { //取字符串中间
	n := strings.Index(str, start)
	if n == -1 {
		return ""
	}
	n += len(start)
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

// StrKillHtml 干掉HTML标记
func StrKillHtml(html string) string {
	r := regexp.MustCompile(`<script[\s\S]*?</script>|<style[\s\S]*?</style>|<[^>]+>|&.{1,8};`)
	return r.ReplaceAllString(html, "")
}

// StrKillBlank 干掉空白和换行
func StrKillBlank(str string) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	return str
}

// HttpsOn 链接Https字符串替换
func HttpsOn(str string) string {
	t, err := beego.AppConfig.Bool("https")
	if err != nil {
		return str
	}
	if t {
		return strings.Replace(str, "http://", "https://", -1)
	}
	return str
}
