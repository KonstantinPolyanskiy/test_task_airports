FROM golang:1.23-alpine AS builder
LABEL authors="polyanskiy"
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/app .
COPY ./stuff /root/stuff
ENTRYPOINT ["./app"]
CMD ["2"]