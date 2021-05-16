FROM golang:alpine AS builder

RUN apk update
WORKDIR /codebase
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /server cmd/server/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /server /server

ENTRYPOINT ["/server"]