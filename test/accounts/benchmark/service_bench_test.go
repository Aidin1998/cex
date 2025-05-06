package benchmark

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"cex/internal/accounts/db"
	"cex/internal/accounts/service"
)

func BenchmarkListAccounts(b *testing.B) {
	// Setup real DB connection once
	dsn := os.Getenv("ACCOUNTS_DSN")
	dbConn, err := db.NewDB(dsn)
	require.NoError(b, err)

	svc := service.NewService(dbConn)
	owner := uuid.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.ListAccounts(context.Background(), owner, 0, 100)
		if err != nil {
			b.Fatal(err)
		}
	}
}
