FROM alpine AS certs

RUN apk add -U --no-cache ca-certificates

FROM golang:alpine AS builder

WORKDIR ./build

COPY cmd/main.go ./
COPY go.mod go.sum ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin

FROM scratch as final

WORKDIR ./app

COPY --from=builder ./build/bin ./
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["./app/bin"]