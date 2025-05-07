package benchmark

import (
	"context"
	"testing"

	"cex/internal/accounts/queue"
	"cex/internal/accounts/service"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func BenchmarkListAccounts(b *testing.B) {
	mockDB, _, _ := sqlmock.New()
	publisher := queue.NewPublisher([]string{"localhost:9092"}, "accounts-events")
	svc := service.NewAccountService(mockDB, publisher)
	ctx := context.Background()
	ownerID := uuid.New()

	for i := 0; i < b.N; i++ {
		_, _ = svc.ListAccounts(ctx, ownerID, 1, 100)
	}
}
