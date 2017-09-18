package runtime

import (
	"fmt"
	"log"
	"net/http"

	"github.com/andy-zhangtao/Canon/util"
	"github.com/julienschmidt/httprouter"
)

const (
	// APIV1 当前API版本
	APIV1 = "/v1"

	// GETVIDEOINFO 获取指定视频实际播放地址
	GETVIDEOINFO = "/video/get"
)

// RunService 运行时态的服务提供商
type RunService struct {
	YtbAPI string
	Port   string
}

// Service 提供RestApi服务
func (r *RunService) Service() error {
	router := httprouter.New()
	router.GET(getAPIPath(""), _testConnect)
	router.GET(getAPIPath(GETVIDEOINFO), r.GetVideoInfo)

	log.Fatal(http.ListenAndServe(":"+r.Port, router))
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
