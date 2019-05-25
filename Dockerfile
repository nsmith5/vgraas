FROM golang:1.12 AS builder
WORKDIR /go/src/github.com/nsmith5/vgraas
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o vgraas cmd/vgraas/*

FROM alpine as production
RUN addgroup -S vgraas && adduser -S vgraas -G vgraas
WORKDIR /home/vgraas
COPY --from=builder /go/src/github.com/nsmith5/vgraas/vgraas .
USER vgraas
CMD ["./vgraas"]
