FROM golang:1.15-alpine

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup \
    && chown appuser:appgroup /app
USER appuser

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o rsc-spreadsheet-api main.go

CMD ["./rsc-spreadsheet-api"]
