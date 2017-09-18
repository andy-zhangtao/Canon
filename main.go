package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/andy-zhangtao/Canon/runtime"
	"github.com/andy-zhangtao/Canon/util"
)

const (
	// RUNTIMETYPE CANON运行状态 更新态还是服务态
	RUNTIMETYPE = "CANON_RUNTIME"
	// RUNPORT CANON运行端口
	RUNPORT = "CANON_RUNTIME_PORT"
	// RUNVIDEOSERVICE 服务态, 负责提供视频API
	RUNVIDEOSERVICE = "CANON_RUNTIME_VIDEO_SERVICE"
	// RUNQUERYSERVICE 服务态, 负责提供视频查询服务
	RUNQUERYSERVICE = "CANON_RUNTIME_QUERY_SERVICE"
	// YTBDAPI youtube-dl APIendpoint
	YTBDAPI = "CANON_YTBD_API"
)

var _VERSION_ = "unknown"

func main() {
	util.SetVersion(_VERSION_)
	fmt.Println(util.GetVersion())

	if ok, name := checkENV(); !ok {
		fmt.Printf("[%s] can not be empty!", name)
		os.Exit(-1)
	}

	var err error
	switch os.Getenv(RUNTIMETYPE) {
	case RUNVIDEOSERVICE:
		port := os.Getenv(RUNPORT)
		if port == "" {
			err = errors.New(RUNPORT + " Can not be Empty!")
			break
		}
		rs := runtime.VideoService{
			Port:   port,
			YtbAPI: os.Getenv(YTBDAPI),
		}
		err = rs.Service()
	case RUNQUERYSERVICE:
		port := os.Getenv(RUNPORT)
		if port == "" {
			err = errors.New(RUNPORT + " Can not be Empty!")
			break
		}
		rs := runtime.QueryService{
			Port: port,
		}
		err = rs.Service()
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}

func checkENV() (bool, string) {
	if os.Getenv(RUNTIMETYPE) == "" {
		return false, RUNTIMETYPE
	}

	if os.Getenv(YTBDAPI) == "" {
		return false, YTBDAPI
	}

	if os.Getenv(RUNPORT) == "" {
		return false, RUNPORT
	}

	return true, ""
}
