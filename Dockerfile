FROM golang:1.12 AS builder
WORKDIR /vgraas
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o vgraas cmd/vgraas/*

FROM alpine as production
RUN addgroup -S vgraas && adduser -S vgraas -G vgraas
WORKDIR /home/vgraas
COPY --from=builder /vgraas/vgraas .
USER vgraas
CMD ["./vgraas"]
