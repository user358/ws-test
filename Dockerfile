FROM golang:1.20

WORKDIR /app
COPY . .

RUN go mod tidy

RUN GO111MODULE=on go install -v ./cmd && go build -o "./bin/ws-test" "./cmd/main.go"

EXPOSE 8080 8081

CMD "/app/bin/ws-test"
