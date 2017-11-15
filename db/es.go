package db

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	vu "github.com/andy-zhangtao/videocrawler/util"

	"gopkg.in/olivere/elastic.v5"
	"github.com/andy-zhangtao/crawlerparam/v1"
)

const (
	// INDEX ES索引名称
	INDEX           = vu.INDEX
	CANON_ES_HOME   = "CANON_ES_HOME"
	CANON_ES_USER   = "CANON_ES_USER"
	CANON_ES_PASSWD = "CANON_ES_PASSWD"
	// EsEmpty ES_HOST变量为空
	EsEmpty = "ES_HOST Can't Be Empty!"
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
	Chanid    []string
	TimeStamp string
	// ID 视频ID
	ID string
}

// Check 检查是否满足DB运行条件
func Check() (bool, string) {
	if os.Getenv(CANON_ES_HOME) == "" {
		return false, CANON_ES_HOME
	}

	if os.Getenv(CANON_ES_USER) == "" {
		return false, CANON_ES_USER
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
		return nil, errors.New(EsEmpty)
	}

	EsUser := os.Getenv("CANON_ES_USER")
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

// getResult 获取Elastic指定条件的数据
func (d *DB) getResult(index, ty string, q elastic.Query) (*elastic.SearchResult, error) {
	return d.client.Search().
		Index(index).
		Type(ty).
		Query(q).
		From(0).
		Size(10).
		Sort("upload", false).
		Pretty(true).
		Do(d.ctx)
}

// GetVideoRangeList 获取指定时间戳之后的视频数据
func (d *DB) GetVideoRangeList() ([]vu.Video, error) {
	var vs []vu.Video

	q := elastic.NewBoolQuery()
	q = q.Filter(elastic.NewRangeQuery("upload").Gt(d.TimeStamp).Lte("now"))
	q = q.Boost(5)

	searchResult, err := d.client.Search().
		Index(d.index).
		Type(d.Ty).
		Query(q).
		From(0).
		Size(10).
		Sort("upload", false).
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
			v.ID = hit.Id
			vs = append(vs, v)
		}
	}
	// for _, cid := range d.Chanid {

	// }

	return vs, nil
}

// GetRandomData 返回随机视频数据
// 用于当前视频物料为空时
func (d *DB) GetRandomData() ([]vu.Video, error) {
	var vs []vu.Video
	q := elastic.NewFunctionScoreQuery().AddScoreFunc(elastic.NewRandomFunction()).Boost(5).MaxBoost(10).BoostMode("multiply")
	searchResult, err := d.client.Search().
		Index(d.index).
		Type(d.Ty).
		Query(q).
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
	// for _, cid := range d.Chanid {
	// 	// termQuery := elastic.NewTermQuery("chanid", cid)

	// }

	return vs, nil
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

// GetInfo 获取指定ID的视频数据
func (d *DB) GetInfo() (vu.Video, error) {
	var vs vu.Video
	termQuery := elastic.NewTermQuery("_id", d.ID)
	searchResult, err := d.client.Search().
		Index(d.index).
	//Type(d.Ty).
		Query(termQuery).
		Pretty(true).
		Do(d.ctx)
	if err != nil {
		return vs, errors.New("Search ElasticSearch By Id Error. " + err.Error())
	}

	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			err = json.Unmarshal(*hit.Source, &vs)
			if err != nil {
				return vs, err
			}
			break
		}
	}

	return vs, nil
}

// GetCZInfo 获取指定ID的锤子视频数据
func (d *DB) GetCZInfo() (vu.CZVideo, error) {
	var vs vu.CZVideo
	termQuery := elastic.NewTermQuery("_id", d.ID)
	searchResult, err := d.client.Search().
		Query(termQuery).
		Pretty(true).
		Do(d.ctx)
	if err != nil {
		return vs, errors.New("Search ElasticSearch By Id Error. " + err.Error())
	}

	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			err = json.Unmarshal(*hit.Source, &vs)
			if err != nil {
				return vs, err
			}
			break
		}
	}

	return vs, nil
}

// GetCZSimilVideo 获取相似视频数据
// chuizi 和 idou 的数据结构不同,必须设置index
func (d *DB) GetCZSimilVideo(index, keys string) ([]vu.CZVideo, error) {
	var vs []vu.CZVideo
	matchQuery := elastic.NewMatchQuery("keys", keys)
	searchResult, err := d.client.Search().
		Index(index).
		Query(matchQuery).
		From(0).
		Size(15).
		Pretty(true).
		Do(d.ctx)
	if err != nil {
		return vs, errors.New("Search ElasticSearch By Id Error. " + err.Error())
	}

	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var v vu.CZVideo
			err = json.Unmarshal(*hit.Source, &v)
			if err != nil {
				return vs, err
			}
			if v.Name == keys {
				continue
			}
			v.ID = hit.Id
			v.URL = "" //URL此时没有必要传过去
			vs = append(vs, v)
		}
	}
	return vs, nil
}

// GetRandomData 返回随机视频数据
// 用于当前视频物料为空时
func (d *DB) GetCZRandomData(index string) ([]interface{}, error) {
	var vs []interface{}
	q := elastic.NewFunctionScoreQuery().AddScoreFunc(elastic.NewRandomFunction()).Boost(5).MaxBoost(10).BoostMode("multiply")
	searchResult, err := d.getResult(index, d.Ty, q)
	if err != nil {
		return vs, errors.New("Search ElasticSearch Error. " + err.Error())
	}

	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var v vu.CZVideo
			err := json.Unmarshal(*hit.Source, &v)
			if err != nil {
				return vs, err
			}
			v.ID = hit.Id
			vs = append(vs, v)
		}
	}

	return vs, nil
}

// GetCZDocInfo 获取指定ID的新闻内容
// index 索引名称
func (d *DB) GetCZDocInfo(index, id string) (interface{}, error) {
	q := elastic.NewTermQuery("_id", id)
	searchResult, err := d.getResult(index, d.Ty, q)
	if err != nil {
		return nil, errors.New("Search ElasticSearch Error. " + err.Error())
	}

	var v v1.Doc
	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			err := json.Unmarshal(*hit.Source, &v)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	return v, nil
}

// GetCZData 获取指定时间戳之后的新闻数据
// index 索引名称
// ty 索引Type名称
// timestamp 时间戳
func (d *DB) GetCZData(index, ty, timestamp string) ([]interface{}, error) {
	var vs []interface{}
	q := elastic.NewBoolQuery()
	q = q.Filter(elastic.NewRangeQuery("upload").Gt(timestamp).Lte("now"))
	q = q.Boost(5)
	searchResult, err := d.getResult(index, ty, q)
	if err != nil {
		return vs, errors.New("Search ElasticSearch Error. " + err.Error())
	}

	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var v v1.Doc
			err := json.Unmarshal(*hit.Source, &v)
			if err != nil {
				return vs, err
			}
			v.ID = hit.Id
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
