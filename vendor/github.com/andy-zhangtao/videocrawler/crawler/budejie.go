package crawler

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/andy-zhangtao/videocrawler/db"
	"github.com/andy-zhangtao/videocrawler/util"
)

const (
	BUDEJIENAME = "budejie"
	BUDEJIEAPI  = "http://m.budejie.com/video/1"
)

type BuDeJie struct {
	URL string
}

// GetChannel 获取百思不得姐的频道列表,因为目前频道较少，直接返回
func (b BuDeJie) GetChannel() (Channel, error) {
	return Channel{
		ChanType: []int{
			0,
		},
	}, nil
}

// ParseChannel 解析频道数据
func (b BuDeJie) ParseChannel(ch int) ([]util.Video, error) {
	var vs []util.Video
	var err error

	switch ch {
	case 0:
		b.URL = BUDEJIEAPI
	}
	doc, err := goquery.NewDocument(b.URL)
	if err != nil {
		return vs, err
	}

	doc.Find(".x-video-p").Each(func(i int, s *goquery.Selection) {
		s.Find("video").Each(func(i int, s *goquery.Selection) {
			var url string
			var title string
			ts, _ := s.Attr("data-tag")
			tts := strings.Split(ts, "|")
			if len(tts) > 3 {
				title = tts[2]
			} else {
				title = ts
			}
			img, _ := s.Attr("poster")
			s.Find("source").Each(func(i int, s *goquery.Selection) {
				url, _ = s.Attr("src")
			})

			vs = append(vs, util.Video{
				Name:     title,
				URL:      url,
				Img:      img,
				Source:   7,
				Desc:     title,
				Keys:     title,
				Upload:   util.GetTimeStamp(),
				ChanID:   ch,
				IsParsed: true,
			})
		})
	})

	return vs, nil
}

// SaveData 数据持久化
func (b BuDeJie) SaveData(vs []util.Video) error {
	db, err := db.GetDB()
	if err != nil {
		return err
	}

	// db.Ty = HUAJIAONAME
	for _, v := range vs {
		db.Ty = util.ChanPara[BUDEJIENAME][v.ChanID]
		if db.Ty == "" {
			log.Println(BUDEJIENAME, v.ChanID)
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
