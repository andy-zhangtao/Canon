FROM alpine:latest

RUN apk update && \
    apk add libc6-compat ca-certificates wget openssl&& \
    update-ca-certificates
    
COPY idou /idou
COPY README.md /readme.md
ADD chan.xml /chan.xml
CMD ["/idou"]