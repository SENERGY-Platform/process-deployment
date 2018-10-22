FROM golang:1.11


COPY . /go/src/process-deployment
WORKDIR /go/src/process-deployment

ENV GO111MODULE=on

RUN go build

EXPOSE 8080

CMD ./process-deployment