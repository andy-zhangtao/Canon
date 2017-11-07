.PHONY: build

# export CRAWLER_URL=http://pupudy.com/xiao
export CRAWLER_NAME=budejie
export CRAWLER_CHAN=0
export CAWLER_DB_TYPE=elasticsearch
export CAWLER_ES_HOME=http://10.50.1.10:9200
export CAWLER_ES_USER=elastic
export CAWLER_ES_PASSWD=changeme
export CRAWLER_RUNTIME_DEBUG=product
export CRAWLER_SERVER_PORT=8081
export CRAWLER_INTERVAL=1

build:
	@go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d)"

run: build
	@./videocrawler

release: *.go util/*.go crawler/*.go db/*.go *.md
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d)" -a -o idou
	docker build -t vikings/idou .
	docker push vikings/idou
