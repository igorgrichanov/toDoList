FROM golang:alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd/
COPY internal internal/
COPY docs docs/
COPY pkg pkg/

RUN go build -o main ./cmd/toDoList
RUN go build -o seed ./cmd/seed

FROM alpine:latest
COPY --from=builder /app/main main
COPY --from=builder /app/seed seed
COPY --from=builder /app/docs docs/
#RUN chmod +x main
CMD ["./main"]