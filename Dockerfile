FROM golang:1.20 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go


FROM gcr.io/distroless/base-debian10
WORKDIR /app
COPY --from=builder /app/main .

ENTRYPOINT ["./main"]