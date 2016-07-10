package controllers

import (
	"github.com/astaxie/beego"
	"github.com/loadfield/go-music/models"
)

type MessageController struct {
	beego.Controller
}

func (c *MessageController) Get() {
	c.Data["Title"] = "留言板"
	c.Data["Name"] = models.NAME
	c.Data["IsMessage"] = true
	c.TplName = "message.html"
}
