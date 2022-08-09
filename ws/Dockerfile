FROM golang:1.17

ENV GO111MODULE="auto"

ENV GOPROXY="https://goproxy.cn"

ARG RT

RUN echo $RT

ENV RUNTIME=$RT

RUN mkdir application

COPY . ./application

WORKDIR "application"

RUN  go build -o main  ./main.go

EXPOSE 3333

CMD ["./main","-f","../etc/megaoasis-api.yaml"]
