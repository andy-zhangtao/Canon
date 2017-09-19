package db

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/andy-zhangtao/bado/util"
	vu "github.com/andy-zhangtao/videocrawler/util"

	elastic "gopkg.in/olivere/elastic.v5"
)

const (
	// INDEX ES索引名称
	INDEX           = vu.INDEX
	CANON_ES_HOME   = "CANON_ES_HOME"
	CANONR_ES_USER  = "CANONR_ES_USER"
	CANON_ES_PASSWD = "CANON_ES_PASSWD"
)

type DB struct {
	// client 数据库具体实例
	client *elastic.Client
	ctx    context.Context
	// Index mapping名称
	index string
	// Ty 类型名称
	Ty string
	// Chanid 频道ID
	Chanid string
}

// Check 检查是否满足DB运行条件
func Check() (bool, string) {
	if os.Getenv(CANON_ES_HOME) == "" {
		return false, CANON_ES_HOME
	}

	if os.Getenv(CANONR_ES_USER) == "" {
		return false, CANONR_ES_USER
	}

	if os.Getenv(CANON_ES_PASSWD) == "" {
		return false, CANON_ES_PASSWD
	}

	return true, ""
}

// GetDB 获取指定名称的数据库实例
func GetDB() (*DB, error) {
	var err error
	db := new(DB)
	db.client, err = returnElastic()
	if err != nil {
		return db, err
	}
	db.ctx = context.Background()
	db.index = INDEX
	return db, nil
}

// returnElastic 返回Elastic实例, 需要BADO_ES_HOME，BADO_ES_USER和BADO_ES_PASSWD
func returnElastic() (*elastic.Client, error) {
	EsHost := os.Getenv("CANON_ES_HOME")
	if EsHost == "" {
		return nil, errors.New(util.EsEmpty)
	}

	EsUser := os.Getenv("CANONR_ES_USER")
	EsPasswd := os.Getenv("CANON_ES_PASSWD")

	// ctx := context.TODO()
	var client *elastic.Client
	var err error

	if (EsUser != "") && (EsPasswd != "") {
		log.Println(EsHost, EsUser, EsPasswd)
		client, err = elastic.NewClient(elastic.SetTraceLog(log.New(os.Stdout, "", 0)), elastic.SetSniff(false), elastic.SetURL(EsHost), elastic.SetBasicAuth(EsUser, EsPasswd))
	} else {
		client, err = elastic.NewClient(elastic.SetTraceLog(log.New(os.Stdout, "", 0)), elastic.SetSniff(false), elastic.SetURL(EsHost))
	}

	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetRandomData 返回随机视频数据
// 用于当前视频物料为空时
func (d *DB) GetRandomData() ([]vu.Video, error) {
	var vs []vu.Video
	termQuery := elastic.NewTermQuery("chanid", d.Chanid)
	q := elastic.NewFunctionScoreQuery().Query(termQuery).AddScoreFunc(elastic.NewRandomFunction()).Boost(5).MaxBoost(10).BoostMode("multiply")
	searchResult, err := d.client.Search().
		Index(d.index).
		Type(d.Ty).
		Query(q).
		// Sort("upload", false).
		From(0).
		Size(10).
		Pretty(true).
		Do(d.ctx)
	if err != nil {
		return vs, errors.New("Search ElasticSearch Error. " + err.Error())
	}

	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var v vu.Video
			err := json.Unmarshal(*hit.Source, &v)
			if err != nil {
				return vs, err
			}

			vs = append(vs, v)
		}
	}

	return vs, nil
	// src, err := q.Source()
	// if err != nil {
	// 	return vs, errors.New("Get Random Video Failed! " + err.Error())
	// }

	// data, err := json.Marshal(src)
	// if err != nil {
	// 	return vs, errors.New("Parse Random Video Failed! " + err.Error())
	// }

	// fmt.Println(string(data))

	// return vs, nil
}

// GetData 获取指定类型的视频数据
func (d *DB) GetData() ([]vu.Video, error) {
	var vs []vu.Video
	termQuery := elastic.NewTermQuery("chanid", d.Chanid)
	searchResult, err := d.client.Search().
		Index(d.index).
		Type(d.Ty).
		Query(termQuery).
		Sort("upload", false).
		From(0).
		Size(10).
		Pretty(true).
		Do(d.ctx)
	if err != nil {
		return vs, errors.New("Search ElasticSearch Error. " + err.Error())
	}

	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var v vu.Video
			err := json.Unmarshal(*hit.Source, &v)
			if err != nil {
				return vs, err
			}

			vs = append(vs, v)
		}
	}

	return vs, nil
}

// // SaveData 保存视频数据到ElasticSearch中
// // vo Video结构体
// func (d *DB) SaveData(vo util.Video) error {

// 	if d.index == "" && d.Ty == "" {
// 		return errors.New("Index or Type is empty!")
// 	}

// 	repeat, err := d.isRepeat(vo.Name)
// 	if err != nil {
// 		return errors.New("Check Repeat Failed! " + err.Error())
// 	}

// 	if repeat {
// 		return errors.New(util.VideoRepeat)
// 	}

// 	return d.save(vo)
// }

// func (d *DB) save(vo util.Video) error {
// 	_, err := d.client.Index().
// 		Index(d.index).
// 		Type(d.Ty).
// 		BodyJson(vo).
// 		Refresh("true").
// 		Do(d.ctx)
// 	if err != nil {
// 		return errors.New("Insert ElasticSearch Error. " + err.Error())
// 	}
// 	return nil
// }

// func (d *DB) isRepeat(title string) (bool, error) {
// 	termQuery := elastic.NewTermQuery("title", title)
// 	searchResult, err := d.client.Search().
// 		Index(d.index).
// 		Query(termQuery).
// 		From(0).
// 		Size(10).
// 		Pretty(true).
// 		Do(d.ctx)
// 	if err != nil {
// 		return true, errors.New("Search ElasticSearch Error. " + err.Error())
// 	}

// 	if searchResult.Hits.TotalHits > 0 {
// 		return true, nil
// 	}

// 	return false, nil
// }
