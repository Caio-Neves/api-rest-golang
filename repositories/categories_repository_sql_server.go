package repositories

import (
	"context"
	"database/sql"
	"rest-api-example/entities"

	sq "github.com/Masterminds/squirrel"
)

type CategoryRepositorySqlServer struct {
	db *sql.DB
}

func NewCategoryRepositorySqlServer(db *sql.DB) *CategoryRepositorySqlServer {
	return &CategoryRepositorySqlServer{
		db: db,
	}
}

func (r CategoryRepositorySqlServer) GetAllCategories(ctx context.Context, params map[string][]string) ([]entities.Category, error) {

	products := sq.Select("id", "name", "description", "active", "createdAt", "updatedAt").From("products")
	if value, exists := params["active"]; exists {
		products = products.Where("active", value[0])
	}
	query, args, err := products.ToSql()
	if err != nil {
		return nil, err
	}
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	var categories []entities.Category
	for rows.Next() {
		var category = entities.Category{}
		err = rows.Scan(&category.Id, &category.Name, &category.Description, &category.Active, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
