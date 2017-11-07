package runtime

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/andy-zhangtao/crawlerparam/v1"
	vu "github.com/andy-zhangtao/videocrawler/util"
	"github.com/julienschmidt/httprouter"
	"github.com/andy-zhangtao/Canon/utils"
	"strings"
	"errors"
)

const (
	CHANIDEMPTY    = "channel id cannot be empty!"
	IDEMPTY        = "video id cannot be empty!"
	INDEXEMPTY     = "index cannot be empty!"
	KEYSEMPTY      = "kyes cannot be empty!"
	TYPESERROR     = "type error!"
	CHANIDNOTEXIST = "this channel does not exist!"
	QUERYVIDEOERR  = "get video list failed!"
	QUERYDOCERROR  = "get doc list failed!"
	PARSEERROR     = "parse object failed!"
	TIMEEMPTY      = "timestamp cannot be empty!"
)

// GetVideoList 获取视频列表
// chanid channel id值
func (q *QueryService) GetVideoList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	channelID := p.ByName("chanid")
	if channelID == "" {
		returnError(errors.New(""), CHANIDEMPTY, w)
		return
	}

	timestamp := p.ByName("time")
	if timestamp == "" {
		returnError(errors.New(""), TIMEEMPTY, w)
		return
	}

	var sc v1.ChanSource
	var vs []vu.Video
	var err error

	// sc = source[util.GetRandom(len(source))]
	q.ESClient.Ty = channelID

	var ncid []string
	for _, i := range sc.CID {
		ncid = append(ncid, strconv.Itoa(i))
	}

	// q.ESClient.Chanid = ncid

	q.ESClient.TimeStamp = timestamp

	vs, err = q.ESClient.GetVideoRangeList()
	if err != nil {
		returnError(err, QUERYVIDEOERR, w)
		return
	}

	returnResult(vs, w)
}

// GetRandomVideoList 获取随机视频列表
// 当前如果没有最新视频物料时，从库存中随机挑选10条视频返回给用户
func (q *QueryService) GetRandomVideoList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	channelID := p.ByName("chanid")
	if channelID == "" {
		returnError(errors.New(""), CHANIDEMPTY, w)
		return
	}

	var sc v1.ChanSource
	var vs []vu.Video
	var err error

	for {
		if len(vs) >= 10 {
			break
		}
		// sc = source[util.GetRandom(len(source))]

		q.ESClient.Ty = channelID
		var ncid []string

		for _, i := range sc.CID {
			ncid = append(ncid, strconv.Itoa(i))
		}

		// q.ESClient.Chanid = ncid

		vs, err = q.ESClient.GetRandomData()
		if err != nil {
			returnError(err, QUERYVIDEOERR, w)
			return
		}
	}

	returnResult(vs, w)
}

// GetVideoInfo 获取视频信息
// id 视频id值
func (q *QueryService) GetVideoInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ID := p.ByName("id")
	if ID == "" {
		returnError(errors.New(""), IDEMPTY, w)
		return
	}

	var vs vu.Video
	var err error

	q.ESClient.ID = ID

	vs, err = q.ESClient.GetInfo()
	if err != nil {
		returnError(err, QUERYVIDEOERR, w)
		return
	}

	returnResult(vs, w)
	return
}

// GetCZVideoInfo 获取视频信息
// id 视频id值
func (q *QueryService) GetCZVideoInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ID := p.ByName("id")
	if ID == "" {
		returnError(errors.New(""), IDEMPTY, w)
		return
	}

	var vs vu.CZVideo
	var err error

	q.ESClient.ID = ID

	vs, err = q.ESClient.GetCZInfo()
	if err != nil {
		returnError(err, QUERYVIDEOERR, w)
		return
	}

	returnResult(vs, w)
}

// GetCZSimilVideoInfo 获取相似视频信息
func (q *QueryService) GetCZSimilVideoInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	qu := r.URL.Query()
	keys := qu.Get("keys")
	if keys == "" {
		returnError(errors.New(""), KEYSEMPTY, w)
		return
	}

	index := qu.Get("index")
	if index == "" {
		returnError(errors.New(""), INDEXEMPTY, w)
		return
	}

	//log.Println(keys)
	var vs []vu.CZVideo
	var err error

	vs, err = q.ESClient.GetCZSimilVideo(index, keys)
	if err != nil {
		returnError(errors.New(""), QUERYVIDEOERR, w)
		return
	}

	returnResult(vs, w)
	return
}

// GetVideoPlayURL 获取视频真实播放地址
func (q *QueryService) GetVideoPlayURL(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	qu := r.URL.Query()
	url := qu.Get("url")
	if url == "" {
		w.WriteHeader(ERROR)
		fmt.Fprintf(w, URLEMPTY)
		return
	}

	if strings.Contains(url, "snssdk.com") {
		//	头条数据
		vs, err := utils.TouTiao(url)
		if err != nil {
			w.WriteHeader(ERROR)
			fmt.Fprintf(w, err.Error())
			return
		}

		respon, err := json.Marshal(vs)
		if err != nil {
			w.WriteHeader(ERROR)
			log.Println(err.Error())
			fmt.Fprintf(w, PARSEERROR)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(respon)
		return
	}

	if strings.Contains(url, "uczzd.cn") {
		vs, err := utils.UC(url)
		if err != nil {
			w.WriteHeader(ERROR)
			fmt.Fprintf(w, err.Error())
			return
		}

		respon, err := json.Marshal(vs)
		if err != nil {
			w.WriteHeader(ERROR)
			log.Println(err.Error())
			fmt.Fprintf(w, PARSEERROR)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(respon)
		return
	}
	w.WriteHeader(ERROR)
	fmt.Fprintf(w, TYPESERROR)
	return
}

// GetRandomVideoList 获取随机视频列表
// 当前如果没有最新视频物料时，从库存中随机挑选10条视频返回给用户
func (q *QueryService) GetCZRandomVideoList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	channelID := p.ByName("chanid")
	if channelID == "" {
		returnError(errors.New(""), CHANIDEMPTY, w)
		return
	}

	var sc v1.ChanSource
	var vs []interface{}
	var err error

	for {
		if len(vs) >= 10 {
			break
		}

		q.ESClient.Ty = channelID
		var ncid []string

		for _, i := range sc.CID {
			ncid = append(ncid, strconv.Itoa(i))
		}

		vs, err = q.ESClient.GetCZRandomData("chuizivideo")
		if err != nil {
			returnError(err, QUERYVIDEOERR, w)
			return
		}
	}

	returnResult(vs, w)
}

func (q *QueryService) GetCZRandomDocList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	channelID := p.ByName("chanid")
	if channelID == "" {
		returnError(errors.New(""), CHANIDEMPTY, w)
		return
	}

	var sc v1.ChanSource
	var vc []interface{}
	var err error

	for {
		if len(vc) >= 10 {
			break
		}

		q.ESClient.Ty = channelID
		var ncid []string

		for _, i := range sc.CID {
			ncid = append(ncid, strconv.Itoa(i))
		}

		vc, err = q.ESClient.GetCZRandomData("chuizidoc")
		if err != nil {
			returnError(err, QUERYDOCERROR, w)
			return
		}
	}
	returnResult(vc, w)
}

func returnResult(o interface{}, w http.ResponseWriter) {
	respon, err := json.Marshal(o)
	if err != nil {
		returnError(err, PARSEERROR, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respon)
	return
}

func returnError(err error, msg string, w http.ResponseWriter) {
	w.WriteHeader(ERROR)
	log.Println(err.Error())
	fmt.Fprintf(w, msg)
}
