package util

import (
	"encoding/xml"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

const (
	// CCTV  cctv
	CCTV = "cctv"
	// WULIU 56视频网
	WULIU = "56.com"
	// BILBIL B站
	BILBIL = "bilbil"
	// DOUYU 斗鱼短视频
	DOUYU = "douyu"
	// HUAJIAO 花椒直播短视频
	HUAJIAO = "huajiao"
	// IQIYI 爱奇艺视频
	IQIYI = "iqiyi"
)

type ChanPara struct {
	XMLName xml.Name   `xml:"chaninfo"`
	Info    []ChanInfo `xml:"chan"`
}

type ChanInfo struct {
	ID     string       `xml:"id,attr"`
	Source []ChanSource `xml:"source"`
}

type ChanSource struct {
	Name string `xml:"name"`
	CID  []int  `xml:"cid"`
}

// type ChanData struct {
// 	Source ChanSource
// }

// MakeChanMap 填充视频列表
func MakeChanMap() (map[string][]ChanSource, error) {
	var cn ChanPara
	fileName := os.Getenv("CANON_CHAN_XML")
	if fileName == "" {
		fileName = "chan.xml"
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &cn)
	if err != nil {
		return nil, err
	}

	cd := make(map[string][]ChanSource)

	for _, cinfo := range cn.Info {
		cd[cinfo.ID] = cinfo.Source
	}

	// log.Println(cd)

	return cd, nil
}

// GetRandom 获取指定范围内的随机数
func GetRandom(length int) int {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rnd.Intn(length)
}
