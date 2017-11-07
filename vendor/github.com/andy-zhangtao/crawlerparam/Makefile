.PHONY: build

export CRAWLER_API_VERSION=v1
export CRAWLER_CALL_BACK=http://127.0.0.1:8081/v1/sync/param


build:
	@go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d)"

run: build
	@./crawlerparam

release: *.go v1/*.go *.xml
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d)" -a -o cralwerparam
	docker build -t vikings/cralwerparam .
	docker push vikings/cralwerparam