# Builder
FROM golang:1.15-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o rsc-spreadsheet-api main.go

# Runner
FROM alpine:3.12

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup \
    && chown appuser:appgroup /app
USER appuser

COPY --chown=appuser:appgroup --from=builder /app/rsc-spreadsheet-api .

CMD ["./rsc-spreadsheet-api"]
