FROM golang:latest
ADD . $GOPATH/src/github.com/Serulian/compiler/.
WORKDIR $GOPATH/src/github.com/Serulian/compiler/.
RUN go get -v ./...
WORKDIR $GOPATH/src/github.com/Serulian/compiler/cmd/serulian/
RUN go build -o /cmd/serulian .
ENTRYPOINT ["/cmd/serulian"]