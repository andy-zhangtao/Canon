.PHONY: build

export CANON_RUNTIME=CANON_RUNTIME_QUERY_SERVICE
export CANON_RUNTIME_PORT=8080
export CANON_YTBD_API=http://10.50.1.10:9191/api/info

export CANON_ES_HOME=http://10.50.1.10:9200
export CANONR_ES_USER=elastic
export CANON_ES_PASSWD=changeme

build:
	@go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d)" -o canon

run: build
	@./canon