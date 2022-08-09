FROM golang:1.17

ENV GO111MODULE="on"

ENV GOPROXY="https://goproxy.cn"

ARG RT

RUN echo $RT

ENV RUNTIME=$RT

RUN mkdir application

COPY . ./application

WORKDIR "application"

RUN  go build -o main ./megaoasis.go

EXPOSE 8888

CMD ["./main","-f","etc/megaoasis-api.yaml"]
