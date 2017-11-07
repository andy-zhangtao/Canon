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
	BIBI         = "bilbil"
	BIBIAPI      = "https://api.bilibili.com/archive_rank/getarchiverankbypartion?tid=%d"
	BIBIVIDEOAPI = "http://www.bilibili.com/video/av%d/"
)

type BILBIL struct {
	URL     string
	Referce string
}
type BIBIVideo struct {
	Data BBVideoList `json:"data"`
}
type BBVideoList struct {
	ViList []BBVideo `json:"archives"`
}
type BBVideo struct {
	ID    int    `json:"aid"`
	Title string `json:"title"`
	Desc  string `json:"description"`
	Img   string `json:"pic"`
}

// GetChannel 获取B站的频道列表,因为目前频道较少，直接返回
func (b BILBIL) GetChannel() (Channel, error) {
	return Channel{
		ChanType: []int{
			21, 75, 76, 138, 160, 161, 162, 163, 174, 175,
		},
	}, nil
}

// ParseChannel 解析频道数据
func (b BILBIL) ParseChannel(ch int) ([]util.Video, error) {
	var vs []util.Video
	var err error
	b.URL = fmt.Sprintf(BIBIAPI, ch)
	log.Println(b.URL)
	client := &http.Client{}
	req, err := http.NewRequest("GET", b.URL, nil)
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

	var bv BIBIVideo

	err = json.Unmarshal(content, &bv)

	if err != nil {
		return nil, err
	}

	for _, bvv := range bv.Data.ViList {
		vs = append(vs, util.Video{
			Name:   bvv.Title,
			Desc:   bvv.Desc,
			URL:    fmt.Sprintf(BIBIVIDEOAPI, bvv.ID),
			Img:    bvv.Img,
			Keys:   bvv.Title,
			Source: 3,
			Upload: util.GetTimeStamp(),
			ChanID: ch,
		})
	}

	return vs, nil
}

// SaveData 数据持久化
func (b BILBIL) SaveData(vs []util.Video) error {
	db, err := db.GetDB()
	if err != nil {
		return err
	}

	// db.Ty = BIBI
	for _, v := range vs {
		db.Ty = util.ChanPara[BIBI][v.ChanID]
		if db.Ty == "" {
			log.Println(BIBI, v.ChanID)
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
