package controllers

import (
	"github.com/astaxie/beego"
	json "github.com/bitly/go-simplejson"
	"github.com/loadfield/go-music/models"
	// "strconv"
	"errors"
	"strings"
)

type XiamiController struct {
	beego.Controller
}

type xiamiSong struct {
	Id          int
	Song_id     int
	Song_name   string
	Album_id    int
	Album_name  string
	Artist_name string
	Artist_logo string
	Album_logo  string
	Listen_file string
}

func (c *XiamiController) Get() {
	id := c.Ctx.Input.Param(":id")
	// beego.Error(id)
	c.Data["Name"] = models.NAME
	c.Data["IsXiami"] = true
	if id != "" {
		song, err := XiamiGet(id)
		if err != nil {
			beego.Error(err)
			c.Redirect("/", 302)
			return
		}
		models.AddSong(models.HomeSong{Url: "/xiami/" + id, Title: song.Song_name, Name: song.Artist_name})
		c.Data["Song"] = song
		c.Data["Title"] = song.Song_name + " | 负荷领域"
		c.TplName = "xiami-music.html"
		return
	}
	c.Data["Title"] = "虾米音乐查询 | 负荷领域"
	c.TplName = "xiami-search.html"
	// xiamiSearch("锦鲤抄")
}

func (c *XiamiController) Post() {
	key := c.Input().Get("key")
	songId := c.Input().Get("songId")
	id := c.Ctx.Input.Param(":id")
	// beego.Error(id, songId, key)
	c.Data["Name"] = models.NAME
	c.Data["IsXiami"] = true
	if id != "" {
		var song xiamiSong
		song.Song_name = c.Input().Get("Song_name")
		song.Album_name = c.Input().Get("Album_name")
		song.Artist_name = c.Input().Get("Artist_name")
		song.Album_logo = c.Input().Get("Album_logo")
		song.Listen_file = c.Input().Get("Listen_file")
		song.Artist_logo = c.Input().Get("Artist_logo")
		// beego.Error(song)
		models.AddSong(models.HomeSong{Url: "/xiami/" + id, Title: song.Song_name, Name: song.Artist_name})
		c.Data["Song"] = song
		c.Data["Title"] = song.Song_name + " | 负荷领域"
		c.TplName = "xiami-music.html"
		return
	} else {
		if key != "" {
			list, err := XiamiSearch(key)
			if err != nil {
				beego.Error(err)
			} else {
				// beego.Error(list)
				c.Data["NotHide"] = true
				c.Data["List"] = list
				c.Data["Value"] = key
			}
		} else if songId != "" {
			songId = models.GetId(songId)
			c.Redirect("/xiami/"+songId, 302)
			return
		}
	}

	c.Data["Title"] = "虾米音乐查询 | 负荷领域"
	c.TplName = "xiami-search.html"
}

func XiamiSearch(key string) (list []xiamiSong, err error) {
	list = make([]xiamiSong, 0)
	str := HttpGet("http://api.xiami.com/web?v=2.0&app_key=1&key="+key+"&page=1&limit=150&r=search/songs", "")
	// beego.Error(string(str))
	j, err := json.NewJson(str)
	if err != nil {
		return
	}
	songs := j.Get("data").Get("songs")
	arr, err := songs.Array()
	if err != nil {
		return
	}
	h := len(arr)
	for i := 0; i < h; i++ {
		var song xiamiSong
		song.Id = i + 1
		song.Song_id, err = songs.GetIndex(i).Get("song_id").Int()
		song.Album_id, err = songs.GetIndex(i).Get("album_id").Int()
		song.Song_name, err = songs.GetIndex(i).Get("song_name").String()
		song.Album_name, err = songs.GetIndex(i).Get("album_name").String()
		song.Artist_name, err = songs.GetIndex(i).Get("artist_name").String()
		song.Artist_logo, err = songs.GetIndex(i).Get("artist_logo").String()
		song.Listen_file, err = songs.GetIndex(i).Get("listen_file").String()
		song.Album_logo, err = songs.GetIndex(i).Get("album_logo").String()
		// beego.Error(song)
		song.Album_logo = strings.Replace(song.Album_logo, "_1.", ".", -1)
		song.Artist_logo = strings.Replace(song.Artist_logo, "_1.", ".", -1)
		if err != nil {
			return
		}
		list = append(list, song)
	}
	// beego.Error(list)
	return
}

func XiamiGet(songId string) (song xiamiSong, err error) {
	str := HttpGet("http://api.xiami.com/web?v=2.0&app_key=1&id="+songId+"&r=song/detail", "")
	j, err := json.NewJson(str)
	if err != nil {
		return
	}
	// beego.Error(string(str))
	song.Song_id, err = j.Get("data").Get("song").Get("song_id").Int()
	song.Album_id, err = j.Get("data").Get("song").Get("album_id").Int()
	song.Song_name, err = j.Get("data").Get("song").Get("song_name").String()
	song.Album_name, err = j.Get("data").Get("song").Get("album_name").String()
	song.Artist_name, err = j.Get("data").Get("song").Get("artist_name").String()
	song.Listen_file, err = j.Get("data").Get("song").Get("listen_file").String()
	song.Album_logo, err = j.Get("data").Get("song").Get("logo").String()
	song.Artist_logo, err = j.Get("data").Get("song").Get("logo").String()
	song.Album_logo = strings.Replace(song.Album_logo, "_1.", ".", -1)
	song.Artist_logo = strings.Replace(song.Artist_logo, "_1.", ".", -1)
	// beego.Error(song)
	if song.Song_id < 1 {
		err = errors.New("虾米音乐.查询失败." + songId)
	}
	return
}
