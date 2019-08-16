FROM golang:1.12 as builder

WORKDIR /app

# Copy the app and build
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build main.go

# Multi-stage binary container from alpine
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/main .
ENTRYPOINT [ "./main" ]
