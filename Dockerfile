FROM golang:1.25

WORKDIR /geo-core


COPY go.mod go.sum ./


RUN go mod download


COPY . .

RUN go build -o /build ./internal/cmd


EXPOSE 8080

CMD ["/build"]