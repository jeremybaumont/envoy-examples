FROM golang:latest as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/rate-limit-mock ./pkg/rate-limit-mock


FROM alpine:latest
COPY --from=builder ./app/bin/rate-limit-mock .
EXPOSE 8081
CMD ["./rate-limit-mock"]
