package database

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository[T any] interface {
	Create(ctx context.Context, entity T) error
	Update(ctx context.Context, entity T) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*T, error)
	FindAll(ctx context.Context, page, pageSize int) ([]T, int64, error)
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type BaseRepository[T any] struct {
	DB *DB
}

func NewBaseRepository[T any](db *DB) *BaseRepository[T] {
	return &BaseRepository[T]{DB: db}
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	return r.DB.WithContext(ctx).Delete(&entity, id).Error
}

func (r *BaseRepository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).First(&entity, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *BaseRepository[T]) FindAll(ctx context.Context, page, pageSize int) ([]T, int64, error) {
	var entities []T
	var total int64

	db := r.DB.WithContext(ctx).Model((*T)(nil))

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *BaseRepository[T]) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
