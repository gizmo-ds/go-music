package controllers

import (
	"github.com/astaxie/beego"
	"github.com/loadfield/go-music/models"
)

type HomeController struct {
	beego.Controller
}

func (c *HomeController) Get() {
	c.Data["SongRecord"] = models.SongRecord
	c.Data["ListRecord"] = models.ListRecord
	c.Data["Title"] = models.NAME + " | 负荷领域"
	c.Data["Name"] = models.NAME
	c.Data["IsHome"] = true
	c.TplName = "home.html"
}
