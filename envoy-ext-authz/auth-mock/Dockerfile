FROM golang:1.17 as builder
WORKDIR /service-proxy
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/auth-mock ./

FROM alpine:latest
WORKDIR /service-proxy/
COPY --from=builder /service-proxy/bin/auth-mock .
EXPOSE 8081
CMD ["./auth-mock"]
