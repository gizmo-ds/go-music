package controllers

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	json "github.com/bitly/go-simplejson"
	"github.com/loadfield/go-music/models"
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
	c.TplName = "netease-song.html"
}

func (c *SongController) Post() {
	songId := c.Input().Get("songId")
	songKey := c.Input().Get("songKey")
	h := c.Input().Get("h") == "on"
	// download := c.Input().Get("download") == "on"
	if songId != "" {
		if strings.Index(songId, "\n") != -1 {
			arr := strings.Split(songId, "\n")
			for i := 0; i < len(arr); i++ {
				arr[i] = models.GetId(arr[i])
			}
			// beego.Error(arr)
			c.Data["Panel"] = true
			List := songList(arr)
			c.Data["H"] = h
			c.Data["List"] = List
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
	} else if songKey != "" {
		List := SongSearch(songKey)
		c.Data["Panel"] = true
		c.Data["List"] = List
		songTpl(c)
	} else {
		c.Redirect("/song", 302)
	}
}

func songList(arr []string) (list []song) {
	for i := 0; i < len(arr); i++ {
		s := GetSongDetail(arr[i])
		s.Id = i + 1
		list = append(list, s)
	}
	return list
}

func GetSongDetail(songId string) (s song) {
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
		s = song{Name: "查询失败", Artists: "查询失败"}
		return
	}
	s.Name, err = j.Get("songs").GetIndex(0).Get("name").String()
	if err != nil {
		s = song{Name: "查询失败", Artists: "查询失败"}
		return
	}
	s.Album, _ = j.Get("songs").GetIndex(0).Get("album").GetIndex(0).Get("name").String()
	s.Artists, _ = j.Get("songs").GetIndex(0).Get("artists").GetIndex(0).Get("name").String()
	s.SongId, _ = strconv.Atoi(songId)
	return
}

func SongSearch(songKey string) (list []song) {
	str := models.HttpPost("http://music.163.com/api/search/pc", "offset=0&limit=100&type=1&s="+songKey)
	j, err := json.NewJson(str)
	if err != nil {
		beego.Error(err)
		return
	}
	code, _ := j.Get("code").Int()
	if code != 200 {
		return
	}
	songs := j.Get("result").Get("songs")
	songCount, _ := j.Get("result").Get("songCount").Int()
	if songCount > 100 {
		songCount = 100
	}
	for i := 0; i < songCount; i++ {
		song := GetSong(songs.GetIndex(i), false)
		song.Id = i + 1
		list = append(list, song)
	}
	return
}

func GetSong(j *json.Json, url bool) (s song) {
	s.SongId, _ = j.Get("id").Int()
	s.Name, _ = j.Get("name").String()
	s.Artists, _ = j.Get("artists").GetIndex(0).Get("name").String()
	s.Album, _ = j.Get("album").Get("name").String()
	s.AlbumId, _ = j.Get("album").Get("id").Int()
	s.AlbumPic, _ = j.Get("album").Get("picUrl").String()
	if url {
		s.Url, _ = j.Get("mp3Url").String()
	}
	return
}
