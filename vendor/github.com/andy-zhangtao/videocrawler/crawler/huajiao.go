package crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/andy-zhangtao/videocrawler/db"
	"github.com/andy-zhangtao/videocrawler/util"
)

const (
	HUAJIAONAME     = "huajiao"
	HUAJIAOAPI      = "http://webh.huajiao.com/live/listVideo?nums=100"
	HUAJIAOVIDEOAPI = "http://www.huajiao.com/v/%d"
)

type HuaJiao struct {
}

type HuaJiaoData struct {
	Data []HuaJiaoVideo `json:"data"`
}
type HuaJiaoVideo struct {
	ID    int    `json:"vid"`
	Title string `json:"video_name"`
	Img   string `json:"video_cover"`
}

// GetChannel 获取花椒短视频的频道列表,因为目前频道较少，直接返回
func (h HuaJiao) GetChannel() (Channel, error) {
	return Channel{
		ChanType: []int{
			0,
		},
	}, nil
}

// ParseChannel 解析频道数据
func (h HuaJiao) ParseChannel(ch int) ([]util.Video, error) {
	var vs []util.Video
	var err error
	// b.URL = fmt.Sprintf(BIBIAPI, ch)
	// log.Println(b.URL)
	client := &http.Client{}
	req, err := http.NewRequest("GET", HUAJIAOAPI, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var hj HuaJiaoData

	err = json.Unmarshal(content, &hj)

	if err != nil {
		return nil, err
	}

	for _, hv := range hj.Data {
		if hv.Title != "" && hv.Title != "null" {
			vs = append(vs, util.Video{
				Name:   hv.Title,
				Desc:   "",
				URL:    fmt.Sprintf(HUAJIAOVIDEOAPI, hv.ID),
				Img:    hv.Img,
				Keys:   hv.Title,
				Source: 5,
				Upload: util.GetTimeStamp(),
				ChanID: ch,
			})
		}
	}

	return vs, nil
}

// SaveData 数据持久化
func (h HuaJiao) SaveData(vs []util.Video) error {
	db, err := db.GetDB()
	if err != nil {
		return err
	}

	// db.Ty = HUAJIAONAME
	for _, v := range vs {
		db.Ty = util.ChanPara[HUAJIAONAME][v.ChanID]
		if db.Ty == "" {
			log.Println(HUAJIAONAME, v.ChanID)
			continue
		}
		err = db.SaveData(v)
		if err != nil {
			if err.Error() == util.VideoRepeat {
				return nil
			}
			return err
		}
	}

	return nil
}
