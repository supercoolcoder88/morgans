FROM golang:1.26.2-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main .

FROM alpine:3.22

RUN apk --no-cache add ca-certificates sqlite-libs tzdata

WORKDIR /app

COPY --from=builder /app/main .

RUN mkdir -p /data

EXPOSE 8000

CMD ["./main"]