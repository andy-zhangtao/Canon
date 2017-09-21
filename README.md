# Canon
更新由videocrawler爬取到的视频url

## How to use?

### ENV List

\* : 必填

env|value|usage|
---|-----|-----|
CANON_RUNTIME|Canon 运行状态|* <br/> CANON_RUNTIME_VIDEO_SERVICE:对外提供视频转换RestAPI <br/> CANON_RUNTIME_QUERY_SERVICE:对外提供视频查询服务|
CANON_RUNTIME_PORT|Canon 监听端口|当CANON_RUNTIME为CANON_RUNTIME_VIDEO_SERVICE时必填|
CANON_YTBD_API|Youtube-dl Endpoint|* |
CANON_CHAN_XML|频道数据配置文件|如果为空，则使用默认的配置数据|

### API List

API Version : v1

Path|Description|
----|-----------|
/v1/|网络测试|

API PATH:

Path|Description|Parameter|
----|-----------|---------|
/video/get|获取指定视频的播放地址|Query Parameter:<br/><br/> url(视频原始页面)|
/video/get/:chanid/:time|获取指定频道的视频列表|URLParameter:<br/><br/>   chanid(频道ID) <br/>  time(时间戳, 13位保留到毫秒)|
/video/random/get/:chanid|获取指定频道的随机列表|URLParameter:<br/><br/> chanid(频道ID)|

### Channel ID List

ID|Description|
--|-----------|
1001|百态视频|
1002|搞笑视频|
1003|发现探索|
1004|体育视频|
1005|音乐配音|
1006|娱乐现场|