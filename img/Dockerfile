FROM golang:1.17

ENV GO111MODULE="auto"

ENV GOPROXY="https://goproxy.cn"

ARG RT

RUN echo $RT

ENV RUNTIME=$RT

RUN mkdir application

COPY . ./application

WORKDIR "application"

#RUN  cd ../

#RUN  go mod tidy

RUN  go build ./filesystem.go

EXPOSE 8090

CMD ["./main"]
