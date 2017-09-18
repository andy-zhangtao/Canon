# Canon
更新由videocrawler爬取到的视频url

## How to use?

### ENV List

\* : 必填

env|value|usage|
---|-----|-----|
CANON_RUNTIME|Canon 运行状态|* <br/> CANON_RUNTIME_VIDEO_SERVICE:对外提供RestAPI <br/> CANON_RUNTIME_QUERY_SERVICE:对外提供视频查询服务|
CANON_RUNTIME_PORT|Canon 监听端口|当CANON_RUNTIME为CANON_RUNTIME_VIDEO_SERVICE时必填|
CANON_YTBD_API|Youtube-dl Endpoint|* |

### API List

API Version : v1

Path|Description|
----|-----------|
/v1/|网络测试|

API PATH:

Path|Description|Parameter|
----|-----------|---------|
/video/get|获取指定视频的播放地址|Query Parameter: url(视频页面)|
/video/get/:chanid|获取指定频道的视频列表|URLParameter: chanid(频道ID)|

### Channel ID List

ID|Description|
1001|社会百态|
1002|社会趣闻|
1003|身边感动|
1004|健康讲堂|
1005|时事新闻|