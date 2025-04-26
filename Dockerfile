FROM golang:1.24

WORKDIR /app

COPY . .

RUN go build -o main ./cmd

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
