FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o web_server ./cmd/main.go

EXPOSE 8080

ENV TODO_PORT="8080"
ENV TODO_DBFILE="scheduler.db"
ENV TODO_PASSWORD="12345"

ENTRYPOINT ["./web_server"]
CMD ["--help"]