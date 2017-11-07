package crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/andy-zhangtao/videocrawler/db"
	"github.com/andy-zhangtao/videocrawler/util"
)

const (
	FISH         = "douyu"
	FISHAPI      = "https://v.douyu.com/video/shortvideo/listAjax?limit=50&page=1&uid=0&type=topic&id=%d"
	FISHVIDEOAPI = "https://v.douyu.com/show/%s"
	FISHREFER    = "https://v.douyu.com/v/s"
)

type DouYu struct {
	URI string
}

type DouYuData struct {
	Data []DouYuVideo `json:"data"`
}

type DouYuVideo struct {
	ID    string `json:"hash_id"`
	Title string `json:"title"`
	Img   string `json:"video_pic"`
	Desc  string `json:"contents"`
}

// GetChannel 获取斗鱼的频道列表,目前斗鱼只有短视频一个频道
func (d DouYu) GetChannel() (Channel, error) {
	return Channel{
		ChanType: []int{
			1, 5, 12, 17, 18, 26, 27, 30,
		},
	}, nil
}

// ParseChannel 解析频道数据
func (d DouYu) ParseChannel(ch int) ([]util.Video, error) {
	var vs []util.Video
	var err error

	d.URI = fmt.Sprintf(FISHAPI, ch)
	client := &http.Client{}
	req, err := http.NewRequest("GET", d.URI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Referer", FISHREFER)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var dyd DouYuData

	err = json.Unmarshal(content, &dyd)
	if err != nil {
		return nil, err
	}

	for _, dy := range dyd.Data {
		if strings.TrimSpace(dy.Title) != "" {
			vs = append(vs, util.Video{
				Name:   dy.Title,
				Desc:   dy.Desc,
				URL:    fmt.Sprintf(FISHVIDEOAPI, dy.ID),
				Img:    dy.Img,
				Keys:   dy.Title,
				Source: 4,
				Upload: util.GetTimeStamp(),
				ChanID: ch,
			})
		}

	}

	return vs, nil
}

// SaveData 数据持久化
func (d DouYu) SaveData(vs []util.Video) error {
	db, err := db.GetDB()
	if err != nil {
		return err
	}

	// db.Ty = FISH
	for _, v := range vs {
		db.Ty = util.ChanPara[FISH][v.ChanID]
		if db.Ty == "" {
			log.Println(FISH, v.ChanID)
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
