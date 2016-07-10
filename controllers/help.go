package controllers

import (
	"github.com/astaxie/beego"
	// json "github.com/bitly/go-simplejson"
	"github.com/loadfield/go-music/models"
	// "strconv"
	// "strings"
)

type HelpController struct {
	beego.Controller
}

func (c *HelpController) Get() {
	c.Data["Title"] = "使用帮助 | 负荷领域"
	c.Data["Name"] = models.NAME
	c.Data["IsHelp"] = true
	c.TplName = "help.html"
}
