FROM alpine:latest

RUN apk update && \
    apk add libc6-compat ca-certificates wget openssl&& \
    update-ca-certificates
    
COPY cralwerparam /cralwerparam
ADD README.md /readme.md
ADD chan.xml /chan.xml
CMD ["/canon"]