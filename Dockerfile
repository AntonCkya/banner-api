FROM golang:1.22

ENV GOPATH=/

COPY . .
RUN go build -o banner-api ./cmd/main.go

RUN go test ./cmd
CMD ["./banner-api"]
