package category

import (
	"context"
	"database/sql"
	"rest-api-example/entities"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CategoryRepositoryPostgres struct {
	db *sql.DB
}

func NewCategoryRepositoryPostgres(db *sql.DB) entities.CategoryInterface {
	return CategoryRepositoryPostgres{
		db: db,
	}
}

func (r CategoryRepositoryPostgres) GetPaginateCategories(ctx context.Context, page int, limit int, params map[string][]string) ([]entities.Category, int, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	countSql := psql.Select("COUNT(*)").
		FromSelect(
			psql.Select("id", "active").
				From("categories"), "subquery",
		)

	if value, exists := params["active"]; exists {
		isActive, err := strconv.Atoi(value[0])
		if err != nil {
			return nil, 0, err
		}
		countSql = countSql.Where("subquery.active = ?", isActive)
	}

	countQuery, countArgs, err := countSql.ToSql()
	if err != nil {
		return nil, 0, err
	}

	// Execute the count query
	var totalCount int
	err = r.db.QueryRow(countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	categoriesSql := psql.Select("id", "name", "description", "active", "created_at", "updated_at").From("categories")
	if value, exists := params["active"]; exists {
		isActive, err := strconv.Atoi(value[0])
		if err != nil {
			return nil, 0, err
		}
		categoriesSql = categoriesSql.Where("active = ?", isActive)
	}
	offset := (page - 1) * limit
	categoriesSql = categoriesSql.Limit(uint64(limit)).Offset(uint64(offset))
	query, args, err := categoriesSql.ToSql()
	if err != nil {
		return nil, 0, err
	}
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, 0, err
	}

	var categories []entities.Category
	for rows.Next() {
		var category = entities.Category{}
		err = rows.Scan(&category.Id, &category.Name, &category.Description, &category.Active, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, category)
	}

	return categories, totalCount, nil
}

func (r CategoryRepositoryPostgres) GetCategoryById(ctx context.Context, id uuid.UUID) (entities.Category, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	categorySql := psql.Select("id", "name", "description", "active", "created_at", "updated_at").From("categories")
	categorySql = categorySql.Where("id = ?", id)

	query, args, err := categorySql.ToSql()
	if err != nil {
		return entities.Category{}, err
	}
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return entities.Category{}, err
	}
	defer stmt.Close()
	row := stmt.QueryRowContext(ctx, args...)
	if row.Err() != nil {
		return entities.Category{}, nil
	}
	category := entities.Category{}
	row.Scan(&category.Id, &category.Name, &category.Description, &category.Active, &category.CreatedAt, &category.UpdatedAt)
	return category, err
}

func (r CategoryRepositoryPostgres) GetCategoriesByIds(ctx context.Context, ids []uuid.UUID) ([]entities.Category, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	categorySql := psql.Select("id", "name", "description", "active", "created_at", "updated_at").From("categories")
	categorySql = categorySql.Where(sq.Eq{"id": ids})

	query, args, err := categorySql.ToSql()
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
		var category entities.Category
		rows.Scan(&category.Id, &category.Name, &category.Description, &category.Active, &category.CreatedAt, &category.UpdatedAt)
		categories = append(categories, category)
	}
	return categories, err
}

func (r CategoryRepositoryPostgres) CreateCategory(ctx context.Context, category entities.Category) (entities.Category, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	categorySql := psql.Insert("categories").Columns("id", "name", "description", "active", "created_at", "updated_at")
	categorySql = categorySql.Values(category.Id, category.Name, category.Description, category.Active, category.CreatedAt, category.UpdatedAt)

	query, args, err := categorySql.ToSql()
	if err != nil {
		return entities.Category{}, err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return entities.Category{}, err
	}
	return category, nil
}

func (r CategoryRepositoryPostgres) DeleteCategoryById(ctx context.Context, id uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	deleteSql := psql.Delete("categories").Where("id = ?", id)
	query, args, err := deleteSql.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r CategoryRepositoryPostgres) DeleteCategories(ctx context.Context, ids []uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	deleteSql := psql.Delete("categories").Where("id = any(?)", pq.Array(ids))
	query, args, err := deleteSql.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r CategoryRepositoryPostgres) UpdateCategoryFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) (entities.Category, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updateSql := psql.Update("categories")
	for key, value := range fields {
		updateSql = updateSql.Set(key, value)
	}
	updateSql = updateSql.Where("id = ?", id)
	query, args, err := updateSql.ToSql()
	if err != nil {
		return entities.Category{}, err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return entities.Category{}, err
	}
	category, err := r.GetCategoryById(ctx, id)
	if err != nil {
		return entities.Category{}, nil
	}
	return category, nil
}

func (r CategoryRepositoryPostgres) GetAllProductsByCategory(ctx context.Context, id uuid.UUID) ([]entities.Product, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	categorySql := psql.Select("p.id", "p.name", "p.description", "p.price", "p.active", "p.created_at", "p.updated_at").From("products_categories")
	categorySql = categorySql.InnerJoin("products p on p.id = products_categories.product_id")
	categorySql = categorySql.Where("category_id = ?", id)

	query, args, err := categorySql.ToSql()
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
	var products []entities.Product
	for rows.Next() {
		var product = entities.Product{}
		err = rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Active, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
