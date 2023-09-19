FROM  golang:1.20

ENV GOPROXY https://goproxy.io,direct
WORKDIR $GOPATH/src/websocket
COPY . $GOPATH/src/websocket
RUN go mod tidy
RUN go build ./cmd/main.go

EXPOSE 3000


