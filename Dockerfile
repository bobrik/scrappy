FROM alpine:3.3

COPY . /go/src/github.com/bobrik/scrappy

RUN apk --update add go git && \
    GOPATH=/go go get github.com/bobrik/scrappy && \
    apk del go git && \
    mv /go/bin/docker-image-cleaner /bin/scrappy && \
    rm -rf /go

ENTRYPOINT ["/bin/scrappy"]
