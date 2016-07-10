package controllers

import (
	"github.com/astaxie/beego"
	json "github.com/bitly/go-simplejson"
	"github.com/loadfield/go-music/models"
	"strconv"
	"strings"
)

type SongController struct {
	beego.Controller
}

func (c *SongController) Get() {
	songTpl(c)
}

func songTpl(c *SongController) {
	c.Data["Title"] = "音乐查询 | 负荷领域"
	c.Data["Name"] = models.NAME
	c.Data["Is163"] = true
	c.TplName = "netease-song.html"
}

func (c *SongController) Post() {
	songId := c.Input().Get("songId")
	h := c.Input().Get("h") == "on"
	if songId != "" {
		if strings.Index(songId, "\n") != -1 {
			arr := strings.Split(songId, "\n")
			for i := 0; i < len(arr); i++ {
				arr[i] = models.GetId(arr[i])
			}
			// beego.Error(arr)
			c.Data["List"] = true
			L := songLise(arr, h)
			c.Data["L"] = L
			c.Data["Value"] = songId
			songTpl(c)
		} else {
			songId = models.GetId(songId)
			if h != false {
				c.Redirect("/music/"+songId+"?h=on", 302)
				return
			}
			c.Redirect("/music/"+songId, 302)
		}
	} else {
		c.Redirect("/song", 302)
	}
}

func songLise(arr []string, h bool) (list []songs) {
	for i := 0; i < len(arr); i++ {
		song := GetSongDetail(arr[i], h)
		song.I = i + 1
		list = append(list, song)
	}
	return list
}

func GetSongDetail(songId string, h bool) (song songs) {
	j, err := json.NewJson(models.HttpGet("http://music.163.com/api/song/detail/?id=" + songId + "&ids=%5B" + songId + "%5D"))
	if err != nil {
		beego.Error("NewJson")
		return
	}
	code, err := j.Get("code").Int()
	if err != nil {
		beego.Error("CodeToInt")
		return
	}
	if code != 200 {
		song = songs{Title: "查询失败", Name: "查询失败"}
		return
	}
	song.Title, err = j.Get("songs").GetIndex(0).Get("name").String()
	if err != nil {
		song = songs{Title: "查询失败", Name: "查询失败"}
		return
	}
	song.Name, _ = j.Get("songs").GetIndex(0).Get("artists").GetIndex(0).Get("name").String()
	song.Id, _ = strconv.Atoi(songId)
	song.H = h
	return
}
