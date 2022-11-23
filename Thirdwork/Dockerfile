FROM golang:alpine 
RUN mkdir -p "/data/app"
WORKDIR "/data/app"
COPY  . . 
RUN pwd
RUN ls .
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN   go build -o server-main server/main.go 
EXPOSE 1235

ENTRYPOINT ["./server-main"]
