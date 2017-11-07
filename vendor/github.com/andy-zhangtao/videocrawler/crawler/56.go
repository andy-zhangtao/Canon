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
	// WULIUNAME 56视频网
	WULIUNAME = "56.com"
	// WULIUAPI 56视频请求API
	WULIUAPI = "http://s1.api.tv.itc.cn/v4/mobile/video/list.json?api_key=%s&column_id=%d&plat=6&column_type=18"
	// WULIUAPIV6 56视频移动端最新视频API
	WULIUAPIV6 = "http://s1.api.tv.itc.cn/v6/mobile/channelPageData/list.json?api_key=%s&plat=6&channel_id=80090000&page_size=20"
	// WULIUAPI = "http://s1.api.tv.itc.cn/v4/mobile/video/list.json?api_key=%s&page_size=100&column_id=%d&column_type=18"
	// WULIUVIDEOAPI 56视频播放页面地址
	WULIUVIDEOAPI = "http://m.56.com/c/v%d.shtml"
)

// WULIU 56视频数据体
type WULIU struct {
	URL string
	Key string
}

// WLData 56视频解析数据体
type WLData struct {
	Data WLDatas `json:"data"`
}

type WLDatas struct {
	Videos []WLVideo `json:"videos"`
}
type WLVideo struct {
	ID    int    `json:"vid"`
	Title string `json:"video_name"`
	Img   string `json:"hor_w8_pic"`
}

// WLDataV6 56视频解析数据体
type WLDataV6 struct {
	Data WLCOLUMNSV6 `json:"data"`
}

type WLCOLUMNSV6 struct {
	Columns []WLDatasV6 `json:"columns"`
}
type WLDatasV6 struct {
	DataList []WUVideoV6 `json:"data_list"`
}
type WUVideoV6 struct {
	ID    int    `json:"vid"`
	Title string `json:"video_name"`
	Img   string `json:"ver_common_pic"`
	Desc  string `json:"tv_desc"`
}

// GetChannel 获取56视频网的频道列表,因为目前频道较少，直接返回
func (w WULIU) GetChannel() (Channel, error) {
	return Channel{
		ChanType: []int{
			0, 615, 616, 617, 618, 619, 1202,
		},
	}, nil
}

// ParseChannel 解析频道数据
func (w WULIU) ParseChannel(ch int) ([]util.Video, error) {
	var vs []util.Video
	var err error
	switch ch {
	case 0:
		w.URL = fmt.Sprintf(WULIUAPIV6, w.Key)
	default:
		w.URL = fmt.Sprintf(WULIUAPI, w.Key, ch)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", w.URL, nil)
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

	switch ch {
	case 0:
		var wl WLDataV6

		err = json.Unmarshal(content, &wl)

		if err != nil {
			return nil, err
		}

		for _, wv := range wl.Data.Columns {
			for _, wvv := range wv.DataList {
				vs = append(vs, util.Video{
					Name:   wvv.Title,
					Desc:   wvv.Desc,
					URL:    fmt.Sprintf(WULIUVIDEOAPI, wvv.ID),
					Img:    wvv.Img,
					Keys:   wvv.Title,
					Source: 2,
					Upload: util.GetTimeStamp(),
					ChanID: ch,
				})
			}
		}
	default:
		var wl WLData

		err = json.Unmarshal(content, &wl)

		if err != nil {
			return nil, err
		}

		for _, wv := range wl.Data.Videos {
			vs = append(vs, util.Video{
				Name:   wv.Title,
				Desc:   "",
				URL:    fmt.Sprintf(WULIUVIDEOAPI, wv.ID),
				Img:    wv.Img,
				Keys:   wv.Title,
				Source: 2,
				Upload: util.GetTimeStamp(),
				ChanID: ch,
			})
		}
	}

	return vs, nil
}

// SaveData 数据持久化
func (w WULIU) SaveData(vs []util.Video) error {
	db, err := db.GetDB()
	if err != nil {
		return err
	}

	// db.Ty = WULIUNAME
	for _, v := range vs {
		db.Ty = util.ChanPara[WULIUNAME][v.ChanID]
		if db.Ty == "" {
			log.Println(WULIUNAME, v.ChanID)
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
