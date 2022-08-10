FROM golang:1.17

ENV GO111MODULE="auto"

ENV GOPROXY="https://goproxy.io"

ARG RT

RUN echo $RT

ENV RUNTIME=$RT

RUN mkdir application

COPY . ./application

WORKDIR "application"

RUN  go build -o main  ./email.go

EXPOSE 3333

CMD ["./main","-f","./etc/email-api.yaml"]
