package crawler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/andy-zhangtao/videocrawler/db"
	"github.com/andy-zhangtao/videocrawler/util"
)

const (
	// CCTVTYPE ES类型名称
	CCTVTYPE = "cctv"
	// CCTVNEWS CCTV新闻
	CCTVNEWS = "http://news.cctv.com/video/data/index.json"
	// CCTVSHEHUI CCTV社会百态频道
	CCTVSHEHUI = "http://news.cctv.com/kuaikan/shehui/data/index.json"
	// CCTVQUWEN CCTV趣闻频道
	CCTVQUWEN = "http://news.cctv.com/kuaikan/quwen/data/index.json"
	// CCTVGANDONG CCTV身边的感动频道
	CCTVGANDONG = "http://news.cctv.com/kuaikan/sbdgd/data/index.json"
)

// CCTV CCTV新闻视频
type CCTV struct {
	URL    string
	ChanID string
}

// CCTVVideo 视频数据
type CCTVVideo struct {
	Title string `json:"title"`
	Desc  string `json:"description"`
	Img   string `json:"image"`
	Keys  string `json:"content"`
	URL   string `json:"url"`
}

type Cvs struct {
	RollData []CCTVVideo `json:"rollData"`
}

// GetChannel 获取CCTV新闻频道，因为CCTV频道较少，直接确定频道。
func (c CCTV) GetChannel() (Channel, error) {
	return Channel{
		ChanType: []int{1001, 1002, 1003, 1004, 1005},
	}, nil
}

// ParseChannel 解析频道数据
func (c CCTV) ParseChannel(ch int) ([]util.Video, error) {
	var vs []util.Video
	switch ch {
	case 1001:
		c.URL = CCTVSHEHUI
	case 1002:
		c.URL = CCTVQUWEN
	case 1003:
		c.URL = CCTVGANDONG
	case 1005:
		c.URL = CCTVNEWS
	default:
		return vs, errors.New("Wrong Channid")
	}

	// c.ChanID = strconv.Itoa(ch)
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.URL, nil)
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

	var cvs Cvs

	err = json.Unmarshal(content, &cvs)

	if err != nil {
		return nil, err
	}

	for _, cv := range cvs.RollData {
		vs = append(vs, util.Video{
			Name:   cv.Title,
			Desc:   cv.Desc,
			URL:    cv.URL,
			Img:    cv.Img,
			Keys:   cv.Keys,
			Source: 0,
			Upload: util.GetTimeStamp(),
			ChanID: ch,
		})
	}

	return vs, nil
}

// SaveData 数据持久化
func (c CCTV) SaveData(vs []util.Video) error {
	db, err := db.GetDB()
	if err != nil {
		return err
	}
	// db.Ty = CCTVTYPE
	for _, v := range vs {
		db.Ty = util.ChanPara[CCTVTYPE][v.ChanID]
		if db.Ty == "" {
			log.Printf(CCTVTYPE, v.ChanID)
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
