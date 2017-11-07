# videocrawler
A crawler for 'eat' videos

## Channel List

ListCode | Name | Source|
---------|------|-------|
10|网络短剧|PuPuDY|
11|搞笑牛人|PuPuDY|
12|搞笑游戏|PuPuDY|
13|搞笑动画|PuPuDY|
14|相声小品|PuPuDY|
15|搞笑综艺|PuPuDY|
16|恶搞配音|PuPuDY|
---------|------|-------|
1001|社会百态|CCTV|
1002|社会趣闻|CCTV|
1003|身边感动|CCTV|
1004|健康讲堂|CCTV|
1005|时事新闻|CCTV|
---------|------|-------|
0|移动端最新视频|56.com|
615|恶搞吐槽|56.com|
616|混剪穿帮|56.com|
617|奇闻趣事|56.com|
618|恶搞配音|56.com|
619|搞笑动画|56.com|
1202|新姿势|56.com|
---------|------|-------|
21|日常|B站|
75|动物圈|B站|
76|美食圈|B站|
138|搞笑|B站|
161|手工|B站|
162|绘画|B站|
163|运动|B站|
174|其它|B站|
175|ASMR|B站|
---------|------|-------|
1|耳朵中毒|斗鱼|
5|卖萌|斗鱼|
12|萌宠|斗鱼|
17|辣眼睛|斗鱼|
18|荣耀时刻|斗鱼|
26|斗鱼好声音|斗鱼|
27|娱乐现场|斗鱼|
30|斗鱼航天局|斗鱼|
---------|------|-------|
0|短视频|花椒
---------|------|-------|
0|欢乐精选|爱奇艺|
1|娱乐八卦|爱奇艺|
2|搞笑短片|爱奇艺|
3|影视吐槽|爱奇艺|
---------|------|-------|
0|短视频|百思不得姐|

## Source List

SourcdID | Name |
---------|------|
0|CCTV|
1|PuPudy|
2|56.com|
3|bilbil|
4|douyu|
5|huajiao|
6|iqiyi|
7|budejie|

## Env List

ENV|Desc|
---|---|
CRAWLER_VIDEO_WULIU_KEY|56视频网请求Key|
CRAWLER_INTERVAL|爬虫运行间隔单位为分|
CRAWLER_RUNTIME_DEBUG|爬虫运行模式 <br/>* debug<br/>* product| 
CRAWLER_SERVER_PORT|爬虫参数同步监听端口|

## Notices

* budejie 视频地址不需要解析