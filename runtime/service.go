package runtime

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/andy-zhangtao/Canon/db"
	"github.com/andy-zhangtao/Canon/util"
	"github.com/julienschmidt/httprouter"
)

const (
	// APIV1 当前API版本
	APIV1 = "/v1"
	// GETVIDEOINFO 获取指定视频实际播放地址
	GETVIDEOINFO = "/video/get"
	// GETVIDEOLIST 获取指定频道的视频数据
	GETVIDEOLIST = "/video/get/:chanid"
	// GETRANDOMVIDEOLIST 获取指定频道的随机视频数据
	GETRANDOMVIDEOLIST = "/video/random/get/:chanid"
)

// VideoService 提供视频地址查询服务
type VideoService struct {
	YtbAPI string
	Port   string
}

// QueryService 提供视频列表查询服务
type QueryService struct {
	Port     string
	ESClient *db.DB
}

// Service 提供RestApi服务
func (v *VideoService) Service() error {
	router := httprouter.New()
	router.GET(getAPIPath(""), _testConnect)
	router.GET(getAPIPath(GETVIDEOINFO), v.GetVideoInfo)

	log.Fatal(http.ListenAndServe(":"+v.Port, router))
	return nil
}

func _testConnect(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(w, "My Name Is LiLei! "+util.GetVersion())
	return
}

// getAPIPath 返回指定版本的API endpoint
func getAPIPath(path string) string {
	log.Println(APIV1 + path)
	return APIV1 + path
}

// Service 提供视频列表
func (q *QueryService) Service() error {
	canRun, name := db.Check()
	if !canRun {
		return errors.New(name + " Cannot Be Empty!")
	}

	var err error
	q.ESClient, err = db.GetDB()
	if err != nil {
		return errors.New("Get ES Client Failed! " + err.Error())
	}

	router := httprouter.New()
	router.GET(getAPIPath(""), _testConnect)
	router.GET(getAPIPath(GETVIDEOLIST), q.GetVideoList)
	router.GET(getAPIPath(GETRANDOMVIDEOLIST), q.GetRandomVideoList)
	log.Fatal(http.ListenAndServe(":"+q.Port, router))
	return nil
}
