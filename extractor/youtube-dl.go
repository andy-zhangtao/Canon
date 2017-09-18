package extractor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// YouTuBeDL 封装后的youtube-dl引擎
type YouTuBeDL struct {
	API string
}

// VideoInfo ytd解析后的视频信息数据
type VideoInfo struct {
	Title     string `json:"title"`
	Desc      string `json:"description"`
	Ext       string `json:"ext"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"url"`
}

// Video ytb返回的视频数据
type Video struct {
	Video []VideoInfo `json:"videos"`
}

// GetVideoInfo 获取指定视频数据
// path 视频页面地址
func (y *YouTuBeDL) GetVideoInfo(path string) (Video, error) {
	var video Video
	requestPath := fmt.Sprintf("%s?url=%s&flatten=true", y.API, path)
	// log.Println(requestPath)
	client := &http.Client{}
	req, err := http.NewRequest("GET", requestPath, nil)
	if err != nil {
		return video, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return video, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return video, err
	}

	// log.Println(string(content))
	err = json.Unmarshal(content, &video)
	if err != nil {
		return video, err
	}

	return video, nil
}
