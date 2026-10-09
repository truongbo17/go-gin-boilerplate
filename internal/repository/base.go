package repository

import (
	"context"
	"errors"

	"github.com/truongbo17/go-gin-boilerplate/internal/page"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaginateParams struct {
	Page         int
	PerPage      int
	Conditions   map[string]any
	ExtraWheres  []WhereClause
	SelectFields []string
	OrderBy      string
	Preloads     []PreloadClause
	Joins        []JoinClause
	GroupBy      string
}

type WhereClause struct {
	Query string
	Args  []any
}

type JoinClause struct {
	Query string
	Args  []any
}

type PreloadClause struct {
	Query string
	Args  []any
}

func NewPaginateParam() *PaginateParams {
	return &PaginateParams{
		Conditions: make(map[string]any),
		Page:       1,
		PerPage:    20,
		OrderBy:    "created_at DESC",
	}
}

type BaseRepository[T any] struct {
	DB *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return BaseRepository[T]{DB: db}
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) Insert(ctx context.Context, entities *[]T) error {
	return r.DB.WithContext(ctx).Create(entities).Error
}

func (r *BaseRepository[T]) InsertIgnore(ctx context.Context, entities *[]T) error {
	return r.DB.WithContext(ctx).
		Clauses(clause.Insert{Modifier: "IGNORE"}).
		Create(entities).
		Error
}

var ErrConflictColumns = errors.New("conflictColumns must not be empty")

func (r *BaseRepository[T]) InsertOrUpdate(ctx context.Context, entity *T, conflictColumns []string, updateColumns []string) error {
	if len(conflictColumns) == 0 {
		return ErrConflictColumns
	}

	return r.DB.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   toClauseColumns(conflictColumns),
			DoUpdates: clause.AssignmentColumns(updateColumns),
		},
	).Create(entity).Error
}

func toClauseColumns(columns []string) []clause.Column {
	clauseColumns := make([]clause.Column, len(columns))
	for i, col := range columns {
		clauseColumns[i] = clause.Column{Name: col}
	}
	return clauseColumns
}

func (r *BaseRepository[T]) WithTransaction(transaction *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{DB: transaction}
}

func (r *BaseRepository[T]) FindByID(ctx context.Context, id uint, entity *T) error {
	return r.DB.WithContext(ctx).First(entity, id).Error
}

func (r *BaseRepository[T]) Save(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T, update any) error {
	return r.DB.WithContext(ctx).Model(entity).Updates(update).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	return r.DB.WithContext(ctx).Delete(&entity, id).Error
}

func (r *BaseRepository[T]) List(ctx context.Context, entities *[]T, conditions map[string]any) error {
	return r.DB.WithContext(ctx).Where(conditions).Find(entities).Error
}

func (r *BaseRepository[T]) FindByCondition(ctx context.Context, condition any, args ...any) ([]T, error) {
	var entities []T
	err := r.DB.WithContext(ctx).Where(condition, args...).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *BaseRepository[T]) FindByConditionForUpdate(ctx context.Context, condition any, args ...any) ([]T, error) {
	var entities []T
	err := r.DB.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where(condition, args...).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *BaseRepository[T]) FindByConditionAndOrder(ctx context.Context, orderBy string, condition any, args ...any) ([]T, error) {
	var entities []T
	query := r.DB.WithContext(ctx).Where(condition, args...)

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	err := query.Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *BaseRepository[T]) FirstByConditionAndOrder(ctx context.Context, orderBy string, condition any, args ...any) (*T, error) {
	var entities T
	query := r.DB.WithContext(ctx).Where(condition, args...)

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	err := query.First(&entities).Error
	if err != nil {
		return nil, err
	}
	return &entities, nil
}

func (r *BaseRepository[T]) FindOneByCondition(ctx context.Context, condition any, args ...any) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).Where(condition, args...).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *BaseRepository[T]) FindOneByConditionAndOrder(ctx context.Context, condition any, order any, args ...any) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).Where(condition, args...).Order(order).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *BaseRepository[T]) FindOneByConditionWithSelect(ctx context.Context, condition, sec any, args ...any) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).Where(condition, args...).Select(sec).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *BaseRepository[T]) DeleteByCondition(ctx context.Context, condition any, args ...any) error {
	var entity T
	result := r.DB.WithContext(ctx).Where(condition, args...).Delete(&entity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *BaseRepository[T]) UpdateByCondition(ctx context.Context, condition any, updates any, args ...any) error {
	var entity *T

	result := r.DB.WithContext(ctx).Model(&entity).Where(condition, args...).Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *BaseRepository[T]) Paginate(ctx context.Context, params *PaginateParams) (*page.Result[T], error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PerPage <= 0 {
		params.PerPage = 20
	}
	if params.PerPage > 100 {
		params.PerPage = 100
	}
	if params.Page > 10000 {
		return nil, errors.New("page exceeds supported range")
	}

	var entities []T
	var total int64

	query := r.DB.WithContext(ctx).Model(&entities)

	if len(params.Conditions) > 0 {
		query = query.Where(params.Conditions)
	}

	for _, c := range params.ExtraWheres {
		query = query.Where(c.Query, c.Args...)
	}

	if len(params.SelectFields) > 0 {
		query = query.Select(params.SelectFields)
	}

	for _, preload := range params.Preloads {
		query = query.Preload(preload.Query, preload.Args...)
	}

	for _, join := range params.Joins {
		query = query.Joins(join.Query, join.Args...)
	}

	if params.GroupBy != "" {
		query = query.Group(params.GroupBy)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PerPage

	if params.OrderBy != "" {
		query = query.Order(params.OrderBy)
	}

	err = query.Offset(offset).Limit(params.PerPage).Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return &page.Result[T]{Items: entities, Number: params.Page, Size: params.PerPage, Total: total}, nil
}
