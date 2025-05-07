# Accounts Service

## Environment Variables
- `ACCOUNTS_PORT`: Port for the service (default: 8081)
- `ACCOUNTS_DSN`: Data source name for the database connection

## Usage
1. Set the required environment variables.
2. Run the service:
   ```
   go run ./cmd/accounts
   ```

## Endpoints
- `GET /healthz`: Health check
- `POST /accounts`: Create a new account
- `GET /accounts/{id}`: Get account details
- `GET /accounts`: List accounts

## Metrics
- Exposed at `/metrics` for Prometheus scraping.