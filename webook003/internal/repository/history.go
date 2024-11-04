package repository

import (
	"context"
	"goworkwebook/webook003/internal/domain"
)

type HistoryRecordRepository interface {
	AddRecord(ctx context.Context, record domain.HistoryRecord) error
}
