package graphql

import (
	"context"
	"strconv"

	core "github.com/truongbo17/go-gin-boilerplate/internal/app/core/graphql"
	"gorm.io/gorm"
)

type MetadataRepository struct{ DB *gorm.DB }

type row struct {
	ID   uint
	Code string
	Name string
}

func (repo MetadataRepository) ListOptions(ctx context.Context, key, keyword string, page, perPage int) ([]core.EntityOption, int64, error) {
	var table, codeColumn, nameColumn string
	switch key {
	case "users":
		table, codeColumn, nameColumn = "users", "username", "username"
	case "roles":
		table, codeColumn, nameColumn = "roles", "slug", "name"
	case "permissions":
		table, codeColumn, nameColumn = "permissions", "slug", "name"
	default:
		return nil, 0, core.ErrUnsupportedKey
	}
	query := repo.DB.WithContext(ctx).Table(table).Where("deleted_at IS NULL")
	if keyword != "" {
		pattern := "%" + keyword + "%"
		query = query.Where("("+codeColumn+" LIKE ? OR "+nameColumn+" LIKE ?)", pattern, pattern)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []row
	if err := query.Select("id", codeColumn+" AS code", nameColumn+" AS name").Order("id ASC").Limit(perPage).Offset((page - 1) * perPage).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	options := make([]core.EntityOption, 0, len(rows))
	for _, item := range rows {
		options = append(options, core.EntityOption{ID: strconv.FormatUint(uint64(item.ID), 10), Code: item.Code, Name: item.Name, Label: item.Name})
	}
	return options, total, nil
}
