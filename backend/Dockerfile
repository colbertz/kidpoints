FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o points-app

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite-libs
WORKDIR /app
COPY --from=builder /app/points-app .
COPY points.db .
EXPOSE 8080
CMD ["./points-app"]
