FROM golang:latest

WORKDIR /app
COPY server server

RUN go build -o server ./server

EXPOSE 8080

CMD ["./server"]
