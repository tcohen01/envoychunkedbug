FROM golang:1.18 as builder

WORKDIR /build

COPY go.mod ./
# COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /build/http_server
RUN chmod 700 /build/http_server

FROM scratch

EXPOSE 8080

WORKDIR /service

COPY --from=builder /build/http_server /service/http_server

CMD ["/service/http_server"]