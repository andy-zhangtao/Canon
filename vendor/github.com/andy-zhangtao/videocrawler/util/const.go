package util

const (
	INTERVAL = "CRAWLER_INTERVAL"
	DEBUG    = "CRAWLER_RUNTIME_DEBUG"
	// INDEX ES索引名称
	INDEX = "idou"
	// SRCYOUKU 优酷视频源
	SRCYOUKU = 0
	// NETSHORTPLAY 网络短剧
	NETSHORTPLAY = "网络短剧"
	// FUNNYCATTLE 搞笑牛人
	FUNNYCATTLE = "搞笑牛人"
	// FUNNYGAME 搞笑游戏
	FUNNYGAME = "搞笑游戏"
	// FUNNYANIMATION 搞笑动画
	FUNNYANIMATION = "搞笑动画"
	// CROSSTALK 相声小品
	CROSSTALK = "相声小品"
	// FUNNYVARIETY 搞笑综艺
	FUNNYVARIETY = "搞笑综艺"
	// FUNNYSPOOF 恶搞配音
	FUNNYSPOOF = "恶搞配音"

	/**
	Elastic Paramers
	*/
	// SaveError 内容保存失败
	SaveError = "Content Save Error"
	// ParseError 内容解析失败
	ParseError = "Parse Error!"
	// DBNameError 不支持的数据库名称
	DBNameError = "Wrong DB Name!"
	// EsEmpty ES_HOST变量为空
	EsEmpty = "ES_HOST Can't Be Empty!"
	// ElasticSearch 数据库名称
	ElasticSearch = "elasticsearch"
	// ElasticType Elastic类型
	ElasticType = 0
	// ElasticTypeError 类型转换失败
	ElasticTypeError = "ElasticSearch Type Error!"
	// VideoRepeat 视频重复
	VideoRepeat = "Video Repeat"
)

// Video 视频数据
type Video struct {
	// ID ES ID标示
	ID string `json:"id"`
	// name 视频名称
	Name string `json:"title"`
	// Desc 视频简介
	Desc string `json:"desc"`
	// URL 视频地址
	URL string `json:"url"`
	// Img 视频缩略图地址
	Img string `json:"img"`
	// Source 视频来源
	Source int `json:"source"`
	// Keys 视频关键字
	Keys string `json:"keys"`
	// Upload 上传时间
	Upload string `json:"upload"`
	// ChanID 频道ID
	ChanID int `json:"chanid"`
	// IsParsed 视频地址是否已经解析过 因为历史原因，当为true时表示已经解析过。 当为false时表示未解析
	IsParsed bool `json:"isparsed"`
}

// CZVideo 锤子视频
type CZVideo struct{
	// ID ES ID标示
	ID string `json:"id"`
	// name 视频名称
	Name string `json:"title"`
	// Desc 视频简介
	Desc string `json:"desc"`
	// URL 视频地址
	URL string `json:"url"`
	// Img 视频缩略图地址
	Img []CZImg `json:"img"`
	// Source 视频来源
	Source int `json:"source"`
	// Keys 视频关键字
	Keys string `json:"keys"`
	// Upload 上传时间
	Upload string `json:"upload"`
	// ChanID 频道ID
	ChanID int `json:"chanid"`
	// IsParsed 视频地址是否已经解析过 因为历史原因，当为true时表示已经解析过。 当为false时表示未解析
	IsParsed bool `json:"isparsed"`
}

// CZImg 锤子视频缩略图
type CZImg struct{
	URL string `json:"url"`
	Width int `json:"width"`
	Height int `json:"height"`
}