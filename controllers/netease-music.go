package controllers

import (
	"crypto/md5"
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	json "github.com/bitly/go-simplejson"
	"github.com/loadfield/go-music/models"
)

type song struct {
	Id       int
	Name     string
	Url      string
	SongId   int
	Album    string
	AlbumId  int
	AlbumPic string
	Artists  string
}

type MusicController struct {
	beego.Controller
}

func (c *MusicController) Get() {
	id := c.Ctx.Input.Param(":id")
	// beego.Error(id)
	H := c.Input().Get("h") == "on"
	if id == "" {
		c.Redirect("/", 302)
		return
	}
	b, s := songDetail(id, H)
	if b == false {
		c.Redirect("/", 302)
		return
	}
	models.AddSong(models.HomeSong{Url: "/music/" + id, Title: s.Name, Name: s.Artists})
	c.Data["Title"] = s.Name + " | 负荷领域"
	c.Data["Song"] = s
	c.Data["Name"] = models.NAME
	c.TplName = "netease-music.html"
}

func songDetail(songId string, h bool) (b bool, song song) {
	s := models.HttpGet("http://music.163.com/api/song/detail/?id=" + songId + "&ids=%5B" + songId + "%5D")
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
		b = false
		return
	}
	song.Name, err = j.Get("songs").GetIndex(0).Get("name").String()
	if err != nil {
		b = false
		return
	}
	song.AlbumId, err = j.Get("songs").GetIndex(0).Get("album").Get("id").Int()
	if !h {
		song.Url, err = j.Get("songs").GetIndex(0).Get("mp3Url").String()
		if song.Url == "" {
			song.Url = albumGetUrl(strconv.Itoa(song.AlbumId), songId, h)
		}
	} else {
		dfsId, _ := j.Get("songs").GetIndex(0).Get("hMusic").Get("dfsId").Int()
		if dfsId == 0 {
			dfsId, _ = j.Get("songs").GetIndex(0).Get("mMusic").Get("dfsId").Int()
		} else if dfsId == 0 {
			dfsId, _ = j.Get("songs").GetIndex(0).Get("lMusic").Get("dfsId").Int()
		}
		encrypted_song_id := encrypt_id(strconv.Itoa(dfsId))
		song.Url = "http://m1.music.126.net/" + encrypted_song_id + "/" + strconv.Itoa(dfsId) + ".mp3"
		if dfsId == 0 {
			song.Url = albumGetUrl(strconv.Itoa(song.AlbumId), songId, h)
		}
	}
	song.Url = strings.Replace(song.Url, "http://m", "http://p", -1)
	song.Album, err = j.Get("songs").GetIndex(0).Get("album").Get("name").String()
	song.AlbumPic, err = j.Get("songs").GetIndex(0).Get("album").Get("picUrl").String()
	song.Artists, err = j.Get("songs").GetIndex(0).Get("artists").GetIndex(0).Get("name").String()
	b = true
	return
}

func albumGetUrl(albumId string, songId string, h bool) (url string) {
	s := models.HttpGet("http://music.163.com/api/album/" + albumId + "?id=" + albumId)
	j, err := json.NewJson(s)
	if err != nil {
		beego.Error("NewJson")
		return ""
	}
	code, err := j.Get("code").Int()
	if err != nil {
		beego.Error("CodeToInt")
		return ""
	}
	if code != 200 {
		return ""
	}
	songs := j.Get("album").Get("songs")
	arr, _ := songs.Array()
	for i := 0; i < len(arr); i++ {
		id, _ := songs.GetIndex(i).Get("id").Int()
		sid, _ := strconv.Atoi(songId)
		if id == sid {
			if !h {
				url, _ = songs.GetIndex(i).Get("mp3Url").String()
			} else {
				dfsId, _ := songs.GetIndex(i).Get("hMusic").Get("dfsId").Int()
				if dfsId == 0 {
					dfsId, _ = songs.GetIndex(i).Get("mMusic").Get("dfsId").Int()
				} else if dfsId == 0 {
					dfsId, _ = songs.GetIndex(i).Get("lMusic").Get("dfsId").Int()
				}
				encrypted_song_id := encrypt_id(strconv.Itoa(dfsId))
				url = "http://m1.music.126.net/" + encrypted_song_id + "/" + strconv.Itoa(dfsId) + ".mp3"
			}
		}
	}
	return
}

// https://github.com/yanunon/NeteaseCloudMusic/wiki/%E7%BD%91%E6%98%93%E4%BA%91%E9%9F%B3%E4%B9%90API%E5%88%86%E6%9E%90
func encrypt_id(id string) string {
	byte1 := []byte("3go8&$8*3*3h0k(2)2")
	byte2 := []byte(id)
	for i := 0; i < len(byte2); i++ {
		byte2[i] = byte2[i] ^ byte1[i%len(byte1)]
	}
	m := md5.New()
	m.Write(byte2)
	s := base64.StdEncoding.EncodeToString(m.Sum(nil))
	m.Reset()
	s = strings.Replace(s, "+", "-", -1)
	s = strings.Replace(s, "/", "_", -1)
	return s
}
