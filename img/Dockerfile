FROM golang:1.17

ENV GO111MODULE="auto"

ENV GOPROXY="https://goproxy.cn"

RUN mkdir application

COPY . ./application

WORKDIR "application"

RUN  go build -o main ./filesystem.go

EXPOSE 8090

CMD ["./main"]
