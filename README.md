# Canon
更新由videocrawler爬取到的视频url

实现的功能如下:

1. 根据视频源实际播放地址
2. 从Idou和Chuizi中获取符合条件的视频数据

## How to use?

### ENV List

\* : 必填

env|value|usage|
---|-----|-----|
CANON_RUNTIME|Canon 运行状态|* <br/> CANON_RUNTIME_VIDEO_SERVICE:对外提供视频转换RestAPI <br/> CANON_RUNTIME_QUERY_SERVICE:对外提供视频查询服务|
CANON_RUNTIME_PORT|Canon 监听端口|当CANON_RUNTIME为CANON_RUNTIME_VIDEO_SERVICE时必填|
CANON_YTBD_API|Youtube-dl Endpoint|* |

### API List

API Version : v1

Path|Description|
----|-----------|
/v1|网络测试|

API PATH:

Path|Description|Parameter|
----|-----------|---------|
/video/get|获取指定视频的播放地址|Query Parameter:<br/><br/> url(视频原始页面)|
/video/get/:chanid/:time|获取指定频道的视频列表|URLParameter:<br/><br/>   chanid(频道ID) <br/>  time(时间戳, 13位保留到毫秒)|
/video/random/get/:chanid|获取指定频道的随机列表|URLParameter:<br/><br/> chanid(频道ID)|
/video/info/:id | 获取指定ID视频信息|URLParameter:<br/><br/> id(视频ID)|
/video/czinfo/:id | 获取指定ID锤子视频信息|URLParameter:<br/><br/> id(视频ID)|
/video/simila | 根据指定关键字查询指定Index中的信息| Query Parameter: <br/><br/> keys(关键字url编码) <br/><br/> index(index名称)|
/video/url | 获取视频真实播放地址 | Query Parameter: <br/><br/> url(原地址)|
/doc/czdata/:chanid/:time|获取指定频道的新闻列表|URLParameter:<br/><br/>   chanid(频道ID) <br/>  time(时间戳, 13位保留到毫秒)|
/doc/czrandom/:chanid|获取指定频道的新闻列表|URLParameter:<br/><br/> chanid(频道ID)|

### Support Platform

1. 头条视频
2. UC视频

### Channel ID List

[点击查看](https://bitbucket.org/andy-zhang/bado/wiki/Channel)