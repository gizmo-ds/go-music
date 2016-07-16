package controllers

import (
	"github.com/astaxie/beego"
	json "github.com/bitly/go-simplejson"
	"github.com/loadfield/go-music/models"
	// "strconv"
)

type AlbunController struct {
	beego.Controller
}

func (c *AlbunController) Get() {
	albumId := c.Ctx.Input.Param(":id")
	albumId = models.GetId(albumId)
	H := c.Input().Get("h") == "on"
	if albumId != "" {
		c.Data["Value"] = albumId
		c.Data["List"], c.Data["Album"] = albumGet(albumId)
		c.Data["NotHide"] = true
		c.Data["H"] = H
	}
	c.Data["Title"] = "专辑查询 | 负荷领域"
	c.Data["Name"] = models.NAME
	c.TplName = "netease-album.html"
}

func (c *AlbunController) Post() {
	albumId := c.Input().Get("albumId")
	albumId = models.GetId(albumId)
	H := c.Input().Get("h") == "on"
	if albumId != "" {
		c.Data["Value"] = albumId
		c.Data["List"], c.Data["Album"] = albumGet(albumId)
		c.Data["NotHide"] = true
		c.Data["H"] = H
	}
	c.Data["Title"] = "专辑查询 | 负荷领域"
	c.Data["Name"] = models.NAME
	c.TplName = "netease-album.html"
}

func albumGet(albumId string) (list []song, title string) {
	s := models.HttpGet("http://music.163.com/api/album/" + albumId + "?id=" + albumId)
	j, err := json.NewJson(s)
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
		return
	}
	title, _ = j.Get("album").Get("name").String()
	songs := j.Get("album").Get("songs")
	arr, _ := songs.Array()
	for i := 0; i < len(arr); i++ {
		var _song song
		_song.SongId, _ = songs.GetIndex(i).Get("id").Int()
		_song.Name, _ = songs.GetIndex(i).Get("name").String()
		_song.Artists, _ = songs.GetIndex(i).Get("artists").GetIndex(0).Get("name").String()
		_song.Id = i + 1
		list = append(list, _song)
		// beego.Error(list[i])
	}
	return
}
