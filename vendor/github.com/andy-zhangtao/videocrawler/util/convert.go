package util

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/andy-zhangtao/crawlerparam/v1"
)

// ChanPara 视频对应码表
// 例如:
// "bilbil":
// 		21:"1001"
var ChanPara map[string]map[int]string

// MakeChanMap 接受同步服务器发来的配置参数，并转换成查询服务器所需要的数据结构
// isLocal 是否使用本地验证。 true 使用本地备份的数据, false 使用同步数据
func MakeChanMap(data []byte, isLocal bool) error {
	var cn v1.ChanPara
	if !isLocal {
		ioutil.WriteFile("chan.xml", data, 0755)
	} else {
		data, _ = ioutil.ReadFile("chan.xml")
	}

	err := xml.Unmarshal(data, &cn)
	if err != nil {
		return err
	}

	if ChanPara == nil {
		ChanPara = make(map[string]map[int]string)
	}

	for _, cinfo := range cn.Info {
		id := cinfo.ID
		for _, s := range cinfo.Source {
			cp := make(map[int]string)
			for _, v := range s.CID {
				cp[v] = id
			}
			tcp := ChanPara[s.Name]
			if tcp != nil {
				for key := range cp {
					tcp[key] = cp[key]
					ChanPara[s.Name] = tcp
				}
			} else {
				ChanPara[s.Name] = cp
			}

		}
	}

	return nil
}

// Convert 将视频名称转换为频道ID
func Convert(name string) int {
	switch name {
	case NETSHORTPLAY:
		return 10
	case FUNNYCATTLE:
		return 11
	case FUNNYGAME:
		return 12
	case FUNNYANIMATION:
		return 13
	case CROSSTALK:
		return 14
	case FUNNYVARIETY:
		return 15
	case FUNNYSPOOF:
		return 16
	}

	return 0
}

// GetTimeStamp 获取采集时间戳
func GetTimeStamp() string {
	return strconv.FormatInt(time.Now().UTC().UnixNano(), 10)[:13]
}

// GetInterVal 返回采集周期
func GetInterVal() int {
	interval := os.Getenv(INTERVAL)
	if interval == "" {
		return 3
	}

	val, err := strconv.Atoi(interval)
	if err != nil {
		return 3
	}

	return val
}

// IsDebug 通过判断CRAWLER_RUNTIME_DEBUG来获取当前运行模式
func IsDebug() bool {
	if strings.ToUpper(os.Getenv(DEBUG)) == "DEBUG" {
		return true
	}
	return false
}
