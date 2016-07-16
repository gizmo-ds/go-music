package controllers

import (
	"github.com/astaxie/beego"
	json "github.com/bitly/go-simplejson"
	"github.com/loadfield/go-music/models"
	// "strconv"
	"crypto/md5"
	"io/ioutil"
	"strings"
	// "crypto/rc4"
	// "encoding/base64"
	"encoding/hex"
	"net/http"
)

type KugouController struct {
	beego.Controller
}

type KugouSong struct {
	Title  string
	Name   string
	Hash   string
	HashMV string
	Album  string
	Url    string
	Pic    string
	H      bool
	Id     int
}

func (c *KugouController) Get() {
	c.Data["IsKugou"] = true
	c.Data["Title"] = "酷狗音乐查询 | 负荷领域"
	c.Data["Name"] = models.NAME
	c.TplName = "kugou-search.html"
}

func (c *KugouController) Post() {
	c.Data["IsKugou"] = true
	keyword := c.Input().Get("key")
	h := c.Input().Get("h") == "on"
	download := c.Input().Get("download") == "on"
	var song KugouSong
	song.Hash = c.Ctx.Input.Param(":hash")
	song.Title = c.Input().Get("Title")
	if song.Hash != "" && song.Title != "" {
		song.Name = c.Input().Get("Name")
		song.Album = c.Input().Get("Album")
		if song.Album == "" {
			song.Album = "NULL"
		}
		song.Url = KugouGetUrl(song.Hash)
		song.Pic = KugouPic(song.Hash)
		c.Data["Song"] = song
		c.Data["Title"] = song.Title + " | 负荷领域"
		c.Data["Name"] = models.NAME
		c.TplName = "kugou-music.html"
		return
	}
	list := KugouSearch(keyword, h, download)
	if download && len(list) <= 30 {
		c.Data["Download"] = true
	}
	c.Data["Value"] = keyword
	c.Data["NotHide"] = true
	c.Data["List"] = list
	c.Data["Title"] = "音乐查询 | 负荷领域"
	c.Data["Name"] = models.NAME
	c.TplName = "kugou-search.html"
}

func KugouPic(hash string) (url string) {
	s := string(HttpGet("http://tools.mobile.kugou.com/api/v1/singer_header/get_by_hash?hash="+hash+"&size=150&format=jsonp", ""))
	j, err := json.NewJson([]byte(s))
	if err != nil {
		beego.Error(err)
		url = "http://models.kugou.com/v3/static/images/index/logo_kugou.png"
		return
	}
	url, err = j.Get("url").String()
	if err != nil {
		url = "http://models.kugou.com/v3/static/images/index/logo_kugou.png"
	}
	return
}

func KugouSearch(keyword string, H, dow bool) (list []KugouSong) {
	keyword = strings.Replace(keyword, " ", "+", -1)
	s := string(HttpGet("http://mobilecdn.kugou.com/api/v3/search/song?format=jsonp&keyword="+keyword+"&page=1&pagesize=50&showtype=1", ""))
	if strings.Index(s, "(") == -1 {
		return
	}
	n := len(s) - 1
	s = string([]byte(s)[1:n])
	j, err := json.NewJson([]byte(s))
	if err != nil {
		beego.Error(err)
		return
	}
	arr, err := j.Get("data").Get("info").Array()
	if err != nil {
		beego.Error(err)
		return
	}
	h := len(arr)
	for i := 0; i < h; i++ {
		var song KugouSong
		info := j.Get("data").Get("info").GetIndex(i)
		song.Title, _ = info.Get("songname").String()
		song.Name, _ = info.Get("singername").String()
		if H {
			song.Hash, err = info.Get("hash320").String()
			if err != nil {
				song.Hash, err = info.Get("hash128").String()
				if err != nil {
					song.Hash, err = info.Get("hash").String()
				}
			}

		} else {
			song.Hash, _ = info.Get("hash").String()
		}
		song.Album, err = info.Get("album_name").String()
		song.Id = i + 1
		song.H = H
		song.HashMV, _ = info.Get("mvhash").String()
		if dow && song.Hash != "" && h <= 30 {
			song.Url = KugouGetUrl(song.Hash)
		}
		list = append(list, song)
	}
	return
}

func HttpGet(url string, cookie string) (body []byte) {
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	reqest.Header.Set("Cookie", cookie)
	reqest.Header.Set("Referer", url)

	response, err := client.Do(reqest)
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

func KugouGetUrl(hash string) (url string) {
	key := hash + "kgcloud"
	m := md5.New()
	m.Write([]byte(key))
	s := HttpGet("http://trackercdn.kugou.com/i/?key="+hex.EncodeToString(m.Sum(nil))+"&cmd=4&acceptMp3=1&pid=1&hash="+hash, "")
	js, err := json.NewJson(s)
	if err != nil {
		beego.Error(err)
		return
	}
	url, err = js.Get("url").String()
	if err != nil {
		beego.Error(err)
		return
	}
	return
}
