package repository

import (
	"context"

	"lumino/internal/domain"
)

type DatasetRepository interface {
	List(ctx context.Context) ([]domain.Dataset, error)
	Get(ctx context.Context, id string) (domain.Dataset, error)
	Create(ctx context.Context, dataset domain.Dataset) error
	Delete(ctx context.Context, id string) error
}
