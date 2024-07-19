# Welcome to Binance Integration!

## To start the project you need to do the following steps:
1. Set up .env, for example
    ```dotenv
   SERVER_PORT=:8000
   API_URL=localhost:8000
   ```
2. Run ``go run cmd/main.go`` to start the project

## After server start on 8000 port
Get all actuall prices from Binance
```curl --location 'http://localhost:8000/api/v1/get-prices```
