package runtime

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/andy-zhangtao/Canon/util"
	vu "github.com/andy-zhangtao/videocrawler/util"
	"github.com/julienschmidt/httprouter"
)

const (
	CHANIDEMPTY    = "channel id cannot be empty!"
	CHANIDNOTEXIST = "this channel does not exist!"
	QUERYVIDEOERR  = "get video list failed!"
	PARSEERROR     = "parse object failed!"
	TIMEEMPTY      = "timestamp cannot be empty!"
)

// GetVideoList 获取视频列表
// chanid channel id值
func (q *QueryService) GetVideoList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	channelID := p.ByName("chanid")
	if channelID == "" {
		w.WriteHeader(ERROR)
		fmt.Fprintf(w, CHANIDEMPTY)
		return
	}

	timestamp := p.ByName("time")
	if timestamp == "" {
		w.WriteHeader(ERROR)
		fmt.Fprintf(w, TIMEEMPTY)
		return
	}

	// chanMap := util.MakeChanMap()
	source := q.ChanMap[channelID]
	if source == nil {
		w.WriteHeader(ERROR)
		fmt.Fprintf(w, CHANIDNOTEXIST)
		return
	}

	var sc util.ChanSource
	var vs []vu.Video
	var err error

	sc = source[util.GetRandom(len(source))]
	q.ESClient.Ty = sc.Name

	var ncid []string
	for _, i := range sc.CID {
		ncid = append(ncid, strconv.Itoa(i))
	}

	q.ESClient.Chanid = ncid

	q.ESClient.TimeStamp = timestamp

	vs, err = q.ESClient.GetVideoRangeList()
	if err != nil {
		w.WriteHeader(ERROR)
		log.Println(err.Error())
		fmt.Fprintf(w, QUERYVIDEOERR)
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
	fmt.Fprintf(w, string(respon))
	return
}

// GetRandomVideoList 获取随机视频列表
// 当前如果没有最新视频物料时，从库存中随机挑选10条视频返回给用户
func (q *QueryService) GetRandomVideoList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	channelID := p.ByName("chanid")
	if channelID == "" {
		w.WriteHeader(ERROR)
		fmt.Fprintf(w, CHANIDEMPTY)
		return
	}

	source := q.ChanMap[channelID]
	if source == nil {
		w.WriteHeader(ERROR)
		fmt.Fprintf(w, CHANIDNOTEXIST)
		return
	}

	var sc util.ChanSource
	var vs []vu.Video
	var err error

	for {
		if len(vs) >= 10 {
			break
		}
		sc = source[util.GetRandom(len(source))]

		q.ESClient.Ty = q.ChanMap[channelID][0].Name
		var ncid []string

		for _, i := range sc.CID {
			ncid = append(ncid, strconv.Itoa(i))
		}

		q.ESClient.Chanid = ncid

		vs, err = q.ESClient.GetRandomData()
		if err != nil {
			w.WriteHeader(ERROR)
			log.Println(err.Error())
			fmt.Fprintf(w, QUERYVIDEOERR)
			return
		}
	}

	respon, err := json.Marshal(vs)
	if err != nil {
		w.WriteHeader(ERROR)
		log.Println(err.Error())
		fmt.Fprintf(w, PARSEERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(respon))
	return
}
