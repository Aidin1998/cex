## Accounts Service

**Usage**  
Start locally:
```bash
ACCOUNTS_PORT=8081 \
DB_URL="postgres://root@localhost:26257/defaultdb?sslmode=disable" \
go run cmd/accounts/main.go
```

**Environment Variables**

- `ACCOUNTS_PORT`: HTTP port for this service (default: 8081)
- `DB_URL`: CockroachDB/Postgres DSN for migrations and runtime

**OpenAPI Spec**  
See `internal/accounts/api/swagger.yaml` for the complete API contract.

**Client SDK**  
If you generated a client, import it:
```go
import "github.com/Aidin1998/cex/clients/accounts"
```