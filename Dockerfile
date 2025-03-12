FROM golang:1.24-bookworm

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.64.5

COPY . .

RUN golangci-lint run

RUN go build -o run-application cmd/main.go

CMD ["./run-application"]
