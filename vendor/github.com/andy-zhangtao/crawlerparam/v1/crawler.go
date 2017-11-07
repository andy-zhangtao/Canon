package v1

import (
	"encoding/xml"
	"io/ioutil"
	"os"
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

	return cd, nil
}

// GetChanMap 获取配置文件内容
func GetChanMap() (string, error) {
	fileName := os.Getenv("CANON_CHAN_XML")
	if fileName == "" {
		fileName = "chan.xml"
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
