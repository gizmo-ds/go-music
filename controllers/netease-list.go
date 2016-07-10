package controllers

import (
	"github.com/astaxie/beego"
	json "github.com/bitly/go-simplejson"
	"github.com/loadfield/go-music/models"
)

type ListController struct {
	beego.Controller
}

type list struct {
	Title       string
	Name        string
	Description string
	List        []songs
}

type songs struct {
	Title string
	Name  string
	Id    int
	I     int
	H     bool
}

func (c *ListController) Get() {
	id := c.Ctx.Input.Param(":id")
	if id != "" {
		listId := id
		H = c.Input().Get("h") != ""
		listId = models.GetId(listId)
		Info := listDetail(listId)
		if len(Info.List) == 0 {
			c.Redirect("/list", 302)
			return
		}
		c.Data["Info"] = Info
		c.Data["NotHide"] = true
	}
	ListTpl(c)
}

var H bool

func (c *ListController) Post() {
	listId := c.Input().Get("listId")
	H = c.Input().Get("h") != ""
	listId = models.GetId(listId)
	Info := listDetail(listId)
	if len(Info.List) == 0 {
		c.Redirect("/list", 302)
		return
	}
	models.AddList(models.HomeList{Id: listId, Views: len(Info.List), Title: Info.Title, Name: Info.Name})
	c.Data["Info"] = Info
	c.Data["Value"] = listId
	c.Data["NotHide"] = true
	ListTpl(c)
}

func ListTpl(c *ListController) {
	c.Data["Title"] = "歌单查询 | 负荷领域"
	c.Data["Name"] = models.NAME
	c.Data["IsList"] = true
	c.TplName = "netease-list.html"
}

func listDetail(listId string) (list list) {
	s := models.HttpGet("http://music.163.com/api/playlist/detail?id=" + listId)
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
	list.Title, _ = j.Get("result").Get("name").String()
	list.Name, _ = j.Get("result").Get("creator").Get("nickname").String()
	list.Description, _ = j.Get("result").Get("description").String()
	l, _ := j.Get("result").Get("tracks").Array()
	h := len(l)
	var song songs
	for i := 0; i < h; i++ {
		song.I = i + 1
		song.Title, _ = j.Get("result").Get("tracks").GetIndex(i).Get("name").String()
		song.Id, _ = j.Get("result").Get("tracks").GetIndex(i).Get("id").Int()
		song.H = H
		song.Name, _ = j.Get("result").Get("tracks").GetIndex(i).Get("artists").GetIndex(0).Get("name").String()
		list.List = append(list.List, song)
	}
	return
}
