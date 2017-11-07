package crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/andy-zhangtao/videocrawler/db"
	"github.com/andy-zhangtao/videocrawler/util"
)

const (
	IQIYINAME = "iqiyi"
	IQIYIAPI  = "http://search.video.iqiyi.com/o?mode=4&ctgName=%s&threeCategory=%s&pageSize=100&type=list&if=html5&site=iqiyi&is_has_father_album=1"
	IQIYICH0  = "欢乐精选"
	IQIYICH1  = "娱乐八卦"
	IQIYICH2  = "搞笑短片"
	IQIYICH3  = "影视剧吐槽"
	IQIYIFLAG = "搞笑"
)

type IQIYI struct {
	URL string
}

type IQIYIData struct {
	Data IQIYIDocInfos `json:"data"`
}
type IQIYIDocInfos struct {
	Infos []IQIYIDoc `json:"docinfos"`
}
type IQIYIDoc struct {
	DocInfo IQIYIDocInfo `json:"albumDocInfo"`
}
type IQIYIDocInfo struct {
	Title string `json:"albumTitle"`
	Desc  string `json:"albumSubTitle"`
	URL   string `json:"albumLink"`
	Img   string `json:"albumHImage"`
}

// GetChannel 获取爱奇艺的频道列表,因为目前频道较少，直接返回
func (i IQIYI) GetChannel() (Channel, error) {
	return Channel{
		ChanType: []int{
			0, 1, 2, 3,
		},
	}, nil
}

// ParseChannel 解析频道数据
func (i IQIYI) ParseChannel(ch int) ([]util.Video, error) {
	var vs []util.Video
	var err error

	switch ch {
	case 0:
		i.URL = fmt.Sprintf(IQIYIAPI, url.QueryEscape(IQIYIFLAG), url.QueryEscape(IQIYICH0))
	case 1:
		i.URL = fmt.Sprintf(IQIYIAPI, url.QueryEscape(IQIYIFLAG), url.QueryEscape(IQIYICH1))
	case 2:
		i.URL = fmt.Sprintf(IQIYIAPI, url.QueryEscape(IQIYIFLAG), url.QueryEscape(IQIYICH2))
	case 3:
		i.URL = fmt.Sprintf(IQIYIAPI, url.QueryEscape(IQIYIFLAG), url.QueryEscape(IQIYICH3))
	}

	log.Println(i.URL)
	client := &http.Client{}
	req, err := http.NewRequest("GET", i.URL, nil)
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

	var iq IQIYIData

	err = json.Unmarshal(content, &iq)

	if err != nil {
		return nil, err
	}

	for _, iqv := range iq.Data.Infos {

		if iqv.DocInfo.Title != "" {
			vs = append(vs, util.Video{
				Name:   iqv.DocInfo.Title,
				Desc:   iqv.DocInfo.Desc,
				URL:    iqv.DocInfo.URL,
				Img:    iqv.DocInfo.Img,
				Keys:   iqv.DocInfo.Title,
				Source: 6,
				Upload: util.GetTimeStamp(),
				ChanID: ch,
			})
		}
	}

	return vs, nil
}

// SaveData 数据持久化
func (i IQIYI) SaveData(vs []util.Video) error {
	db, err := db.GetDB()
	if err != nil {
		return err
	}

	// db.Ty = IQIYINAME
	for _, v := range vs {
		db.Ty = util.ChanPara[IQIYINAME][v.ChanID]
		if db.Ty == "" {
			log.Println(IQIYINAME, v.ChanID)
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
