package db

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/andy-zhangtao/videocrawler/util"

	elastic "gopkg.in/olivere/elastic.v5"
)

type DB struct {
	// client 数据库具体实例
	client *elastic.Client
	ctx    context.Context
	// Index mapping名称
	index string
	// Ty 类型名称
	Ty string
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
	db.index = util.INDEX
	return db, nil
}

// Save 保存数据到指定数据库中
// func (d *DB) Save() error {
// 	return nil
// }

// returnElastic 返回Elastic实例, 需要BADO_ES_HOME，BADO_ES_USER和BADO_ES_PASSWD
func returnElastic() (*elastic.Client, error) {
	EsHost := os.Getenv("CAWLER_ES_HOME")
	if EsHost == "" {
		return nil, errors.New(util.EsEmpty)
	}

	EsUser := os.Getenv("CAWLER_ES_USER")
	EsPasswd := os.Getenv("CAWLER_ES_PASSWD")

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

// SaveData 保存视频数据到ElasticSearch中
// vo Video结构体
func (d *DB) SaveData(vo util.Video) error {

	if d.index == "" && d.Ty == "" {
		return errors.New("Index or Type is empty!")
	}

	repeat, err := d.isRepeat(vo.Name)
	if err != nil {
		return errors.New("Check Repeat Failed! " + err.Error())
	}

	if repeat {
		return errors.New(util.VideoRepeat)
	}

	return d.save(vo)
}

func (d *DB) save(vo util.Video) error {
	_, err := d.client.Index().
		Index(d.index).
		Type(d.Ty).
		BodyJson(vo).
		Refresh("true").
		Do(d.ctx)
	if err != nil {
		return errors.New("Insert ElasticSearch Error. " + err.Error())
	}
	return nil
}

func (d *DB) isRepeat(title string) (bool, error) {
	termQuery := elastic.NewTermQuery("title", title)
	searchResult, err := d.client.Search().
		Index(d.index).
		Query(termQuery).
		From(0).
		Size(10).
		Pretty(true).
		Do(d.ctx)
	if err != nil {
		return true, errors.New("Search ElasticSearch Error. " + err.Error())
	}

	if searchResult.Hits.TotalHits > 0 {
		return true, nil
	}

	return false, nil
}
