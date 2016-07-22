package routers

import (
	"github.com/astaxie/beego"
	"github.com/loadfield/go-music/controllers"
)

func init() {
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/album", &controllers.AlbunController{})
	beego.Router("/album/:id", &controllers.AlbunController{})
	beego.Router("/song", &controllers.SongController{})
	beego.Router("/xiami", &controllers.XiamiController{})
	beego.Router("/xiami/:id", &controllers.XiamiController{})
	beego.Router("/kugou", &controllers.KugouController{})
	beego.Router("/kugou/:hash", &controllers.KugouController{})
	beego.Router("/list", &controllers.ListController{})
	beego.Router("/list/:id", &controllers.ListController{})
	beego.Router("/message", &controllers.MessageController{})
	beego.Router("/music", &controllers.MusicController{})
	beego.Router("/music/:id", &controllers.MusicController{})
}
