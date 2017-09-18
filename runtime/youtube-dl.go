package runtime

import (
	"fmt"
	"log"
	"net/http"

	"github.com/andy-zhangtao/Canon/extractor"
	"github.com/julienschmidt/httprouter"
)

const (
	// ERROR 通用错误码
	ERROR = 501
	// URLEMPTY URL参数为空
	URLEMPTY = "URL Can not be empty!"
	// PARSE_VIDEO_ERROR 解析视频数据失败
	PARSE_VIDEO_ERROR = "Get Video Failed!"
	// VIDEO_LENGTH_ERROR 视频长度错误
	VIDEO_LENGTH_ERROR = "Video Length Error!"
)

// GetVideoInfo 获取视频实际播放地址
func (rs *VideoService) GetVideoInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	q := r.URL.Query()
	url := q.Get("url")

	if url == "" {
		w.WriteHeader(ERROR)
		fmt.Fprintf(w, URLEMPTY)
		return
	}

	log.Printf("Request path:[%s] \n", url)
	ytbd := extractor.YouTuBeDL{
		API: rs.YtbAPI,
	}

	video, err := ytbd.GetVideoInfo(url)
	if err != nil {
		w.WriteHeader(ERROR)
		log.Println(err.Error())
		fmt.Fprintf(w, PARSE_VIDEO_ERROR)
		return
	}

	if len(video.Video) == 0 {
		w.WriteHeader(ERROR)
		fmt.Fprintf(w, VIDEO_LENGTH_ERROR)
		return
	}

	fmt.Fprintf(w, video.Video[0].URL)
	return
}
