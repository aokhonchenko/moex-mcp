FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o moex-mcp .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/moex-mcp /usr/local/bin/moex-mcp

ENTRYPOINT ["moex-mcp"]
