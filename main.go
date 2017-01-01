package main

import (
	"github.com/astaxie/beego"
	"github.com/loadfield/go-music/models"
	_ "github.com/loadfield/go-music/routers"
)

func main() {
	beego.AddFuncMap("https", models.HttpsOn)
	beego.Run()
}
