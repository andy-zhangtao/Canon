package crawler

import (
	"errors"
	"log"
	"time"

	"github.com/andy-zhangtao/videocrawler/util"
)

// Channel 媒体频道
type Channel struct {
	// ChanType 频道类型
	ChanType []int
}

// Crawler 爬虫接口定义
type Crawler interface {
	GetChannel() (Channel, error)
	ParseChannel(ch int) ([]util.Video, error)
	SaveData(vs []util.Video) error
}

// Do 启动爬虫
// c 爬虫类型实例
// chs 频道类型
func Do(c Crawler, chs map[int]int) error {
	ch, err := c.GetChannel()
	if err != nil {
		return errors.New("GetChannel Failed! " + err.Error())
	}

	if util.IsDebug() {
		for _, channel := range ch.ChanType {
			var vs []util.Video
			log.Println(channel)
			if chs[channel] == 1 {
				tvs, err := c.ParseChannel(channel)
				if err != nil {
					return errors.New("Parse Channel Failed! " + err.Error())
				}

				vs = append(vs, tvs...)
			}
			err = c.SaveData(vs)
			if err != nil {
				return errors.New("Save Data Failed! " + err.Error())
			}
		}

	} else {
		for {
			now := time.Now()
			next := now.Add(time.Minute * time.Duration(util.GetInterVal()))
			next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())
			ti := time.NewTimer(next.Sub(now))
			log.Printf("下次采集时间为[%s]\n", next.Format("200601021504"))
			select {
			case <-ti.C:
				for _, channel := range ch.ChanType {
					var vs []util.Video
					// log.Println(channel)
					if chs[channel] == 1 {
						tvs, err := c.ParseChannel(channel)
						if err != nil {
							return errors.New("Parse Channel Failed! " + err.Error())
						}

						vs = append(vs, tvs...)
					}
					err = c.SaveData(vs)
					if err != nil {
						return errors.New("Save Data Failed! " + err.Error())
					}
				}
			}
		}
	}

	return nil
}
