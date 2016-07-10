package main

import (
	"github.com/astaxie/beego"
	_ "github.com/loadfield/go-music/routers"
)

func main() {
	beego.Run()
}
