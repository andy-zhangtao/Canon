FROM alpine:latest

RUN apk update && \
    apk add libc6-compat ca-certificates wget openssl&& \
    update-ca-certificates
    
COPY canon /canon
ADD README.md /readme.md
CMD ["/canon"]