# ws-test

Simple service for testing various protocols such as http, grpc, websockets.

#### Tests
```Bash
go test ./...
```

#### Set config by environment variables:
```Bash
export WS_TEST_HTTP_PORT=8080
export WS_TEST_GRPC_PORT=8081
```

#### Run
```Bash
go run ./cmd/main.go
```

#### Example requests

Get user balance:
```Bash
curl --location --request GET 'localhost:8080/api/wallet/balance/1'
```

Make a deposit:
```Bash
curl --location --request POST 'http://localhost:8080/api/wallet/deposit' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id": 1,
    "value": "200"
}'
```

Withdraw funds:
```Bash
curl --location --request POST 'http://localhost:8080/api/wallet/withdraw' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id": 1,
    "value": "100"
}'
```

#### Test websocket client
```Bash
go run ./cmd/ws_client/ws_client.go -addr localhost:8080
```


#### Docker Example
```Bash
docker build -t ws-test .
docker run -d -p 8080:8080 -p 8081:8081 ws-test 
```