FROM golang:1.16.2-buster AS builder
WORKDIR /opt/golangbb/builddir
COPY . .
RUN go build ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /opt/golangbb
COPY --from=builder /opt/golangbb/builddir/main .
CMD ["./main"]