FROM golang:1.15.3
WORKDIR $GOPATH/src/github.com/silverswords/muses.minio/main
COPY . $GOPATH/src/github.com/silverswords/muses.minio/
RUN go build -o minio-api ./main.go
EXPOSE 8000
ENTRYPOINT ["./minio-api"]