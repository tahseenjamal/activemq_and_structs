FROM golang:1.21.5

WORKDIR /app

EXPOSE 61613

COPY . .

RUN go build main.go

ENTRYPOINT ["/app/main"]
