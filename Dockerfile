FROM golang:1.25

WORKDIR /geo-notification-system
COPY . .

RUN go build -o /build ./internal/cmd


EXPOSE 8080

CMD ["/build"]