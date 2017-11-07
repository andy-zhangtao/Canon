package crawler

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/andy-zhangtao/videocrawler/util"
)

const (
	// PFXURL 地址前缀
	PFXURL = "http://pupudy.com"
)

// PuPuDy 噗噗电影网
type PuPuDy struct {
	URL string
}

// GetChannel 获取频道列表
func (p PuPuDy) GetChannel() (Channel, error) {

	channel := Channel{}

	doc, err := goquery.NewDocument(p.URL)
	if err != nil {
		return channel, err
	}

	doc.Find(".js-tongjip").Each(func(i int, s *goquery.Selection) {
		channel.ChanType = append(channel.ChanType, util.Convert(s.Text()))
	})

	return channel, nil
}

// ParseChannel 解析频道数据
func (p PuPuDy) ParseChannel(ch int) ([]util.Video, error) {
	switch ch {
	case 10:
		p.URL += "?make=all"
	case 11:
		p.URL += "?make=gaoxiaoniuren"
	case 12:
		p.URL += "?make=youxi"
	case 13:
		p.URL += "?make=dongman"
	case 14:
		p.URL += "?make=xiaopinxiangsheng"
	case 15:
		p.URL += "?make=zongyi"
	case 16:
		p.URL += "?make=peiyin"
	}

	var vs []util.Video

	doc, err := goquery.NewDocument(p.URL)
	if err != nil {
		return vs, err
	}

	doc.Find(".mod-pic").Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			alt, _ := s.Find("img").Attr("alt")
			src, _ := s.Find("img").Attr("src")

			if len(alt) > 0 && len(src) > 0 {

				tempUrl := PFXURL + href
				log.Printf("Parse [%s] \n", tempUrl)
				tdoc, err := goquery.NewDocument(tempUrl)
				if err != nil {
					log.Printf("Parse [%s] Failed! [%s]\n", tempUrl, err.Error())
				} else {
					if turl, ok := tdoc.Find("iframe").Attr("src"); ok {
						if strings.Contains(turl, "?url=") {
							us := strings.Split(turl, "?url=")
							if len(us) > 1 {
								vs = append(vs, util.Video{
									Name:   alt,
									URL:    us[1],
									Img:    src,
									Source: 1,
								})
							}
						}
					}
				}

			}
		})
	})
	return vs, nil
}

// SaveData 数据持久化
func (p PuPuDy) SaveData(vs []util.Video) error {

	// db, err := db.GetDB(util.ElasticSearch)
	// if err != nil {
	// 	panic(err)
	// }

	return nil
}
