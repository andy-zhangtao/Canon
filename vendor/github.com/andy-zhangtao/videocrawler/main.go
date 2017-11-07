package main

import (
	"fmt"
	"httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/andy-zhangtao/videocrawler/crawler"
	"github.com/andy-zhangtao/videocrawler/util"
)

const (
	ERROR      = 501
	SERVERPORT = "CRAWLER_SERVER_PORT"
	// CRAWLER_URL 目标URL
	CRAWLER_URL = "CRAWLER_URL"
	// DESTNAME 目标网址名称,用来区别爬虫规则
	DESTNAME = "CRAWLER_NAME"
	// DESTCHAN 需要解析的频道
	DESTCHAN = "CRAWLER_CHAN"
	// DBTYPE 数据库类型
	DBTYPE = "CAWLER_DB_TYPE"
	// ESHOME es地址
	ESHOME = "CAWLER_ES_HOME"
	// ESUSER es用户名
	ESUSER = "CAWLER_ES_USER"
	// ESPASSWD es口令
	ESPASSWD = "CAWLER_ES_PASSWD"
	// PUPUDY 噗噗电影网
	PUPUDY = "pupudy.com"
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
	// BUDEJIE 百思不得姐
	BUDEJIE = "budejie"
	// WULIUKEY 56视频请求KEY
	WULIUKEY = "CRAWLER_VIDEO_WULIU_KEY"
)

var _VERSION_ = "unknown"

func main() {

	fmt.Printf("============================\n")
	fmt.Printf("Video Crawler Version [%s]\n", _VERSION_)
	fmt.Printf("============================\n")
	if isOK, env := checkENV(); !isOK {
		fmt.Printf("ENV [%s] check failed! \n", env)
		os.Exit(-1)
	}

	go syncpara()

	if len(util.ChanPara) == 0 {
		util.MakeChanMap(nil, true)
	}

	var err error
	chm := getChans()
	name := os.Getenv(DESTNAME)
	switch name {
	case PUPUDY:
		if os.Getenv(CRAWLER_URL) == "" {
			fmt.Println("Should Specify Env CRAWLER_URL")
			os.Exit(-1)
		}
		pu := crawler.PuPuDy{
			URL: os.Getenv(CRAWLER_URL),
		}
		err = crawler.Do(pu, chm)

	case CCTV:
		c := crawler.CCTV{}
		err = crawler.Do(c, chm)
	case WULIU:
		key := os.Getenv(WULIUKEY)
		if key == "" {
			fmt.Println("Please Specify Env CRAWLER_VIDEO_WULIU_KEY, if you want get data from 56.com")
			os.Exit(-1)
		}

		w := crawler.WULIU{
			Key: key,
		}

		err = crawler.Do(w, chm)
	case BILBIL:
		b := crawler.BILBIL{}
		err = crawler.Do(b, chm)

	case DOUYU:
		d := crawler.DouYu{}
		err = crawler.Do(d, chm)
	case HUAJIAO:
		h := crawler.HuaJiao{}
		err = crawler.Do(h, chm)
	case IQIYI:
		i := crawler.IQIYI{}
		err = crawler.Do(i, chm)
	case BUDEJIE:
		b := crawler.BuDeJie{}
		err = crawler.Do(b, chm)
	default:
		err = fmt.Errorf("The [%s] maybe not in the configure. ", name)
	}

	if err != nil {
		fmt.Printf("Error [%s] \n", err.Error())
		os.Exit(-1)
	}

	return
}

func checkENV() (bool, string) {
	if os.Getenv(util.DEBUG) == "" {
		return false, util.DEBUG
	}

	if os.Getenv(DESTNAME) == "" {
		return false, DESTNAME
	}

	if os.Getenv(DESTCHAN) == "" {
		return false, DESTCHAN
	}

	if os.Getenv(DBTYPE) == "" {
		return false, DBTYPE
	}
	if os.Getenv(ESHOME) == "" {
		return false, ESHOME
	}

	if os.Getenv(ESUSER) == "" {
		return false, ESUSER
	}

	if os.Getenv(ESPASSWD) == "" {
		return false, ESPASSWD
	}

	if os.Getenv(SERVERPORT) == "" {
		return false, SERVERPORT
	}
	return true, ""
}

func getChans() map[int]int {
	chs := os.Getenv(DESTCHAN)

	cs := strings.Split(chs, ",")

	chm := make(map[int]int)

	for _, cc := range cs {
		tc, _ := strconv.Atoi(cc)
		chm[tc] = 1
	}

	return chm
}

func syncpara() {
	router := httprouter.New()
	router.POST("/v1/sync/param", _syncpara)
	log.Fatal(http.ListenAndServe(":"+os.Getenv(SERVERPORT), router))
}

func _syncpara(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(ERROR)
		log.Println(err.Error())
		return
	}

	err = util.MakeChanMap(data, false)
	if err != nil {
		w.WriteHeader(ERROR)
		log.Println(err.Error())
		return
	}

	log.Println("接受到最新配置文件")
	return
}
