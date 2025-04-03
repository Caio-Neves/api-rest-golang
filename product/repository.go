package product

import (
	"context"
	"database/sql"
	"rest-api-example/entities"
	"strconv"

	"github.com/lib/pq"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type ProductRepositoryPostgres struct {
	db *sql.DB
}

func NewProductRepositoryPostgres(db *sql.DB) entities.ProductInterface {
	return ProductRepositoryPostgres{
		db: db,
	}
}

func (r ProductRepositoryPostgres) GetAllProducts(ctx context.Context, filters map[string][]string) ([]entities.Product, int, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	countSql := psql.Select("COUNT(*)").
		FromSelect(
			psql.Select("id", "active").
				From("categories"), "subquery",
		)

	if value, exists := filters["active"]; exists {
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

	productSql := psql.Select("id", "name", "description", "price", "active", "created_at", "updated_at").From("products")
	if value, exists := filters["active"]; exists {
		isActive, err := strconv.Atoi(value[0])
		if err != nil {
			return nil, 0, err
		}
		productSql = productSql.Where("active = ?", isActive)
	}
	page := 1
	limit := 10
	if value, exists := filters["page"]; exists {
		page, _ = strconv.Atoi(value[0])
	}
	if value, exists := filters["limit"]; exists {
		limit, _ = strconv.Atoi(value[0])
	}
	offset := (page - 1) * limit
	productSql = productSql.Limit(uint64(limit)).Offset(uint64(offset))

	query, args, err := productSql.ToSql()
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

	var products []entities.Product
	for rows.Next() {
		var product = entities.Product{}
		err = rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Active, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}
	return products, totalCount, nil
}

func (r ProductRepositoryPostgres) GetProductById(ctx context.Context, id uuid.UUID) (entities.Product, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	productSql := psql.Select("id", "name", "description", "price", "active", "created_at", "updated_at").From("products")
	productSql = productSql.Where("id = ?", id)

	query, args, err := productSql.ToSql()
	if err != nil {
		return entities.Product{}, err
	}
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return entities.Product{}, err
	}
	defer stmt.Close()
	row := stmt.QueryRowContext(ctx, args...)
	if row.Err() != nil {
		return entities.Product{}, nil
	}
	product := entities.Product{}
	row.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Active, &product.CreatedAt, &product.UpdatedAt)
	return product, err
}

func (r ProductRepositoryPostgres) DeleteProductById(ctx context.Context, id uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	deleteSql := psql.Delete("products").Where("id = ?", id)
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

func (r ProductRepositoryPostgres) DeleteProducts(ctx context.Context, ids []uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	deleteSql := psql.Delete("products").Where("id = any(?)", pq.Array(ids))
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

func (r ProductRepositoryPostgres) CreateProduct(ctx context.Context, product entities.Product) (entities.Product, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	productSql := psql.Insert("products").Columns("id", "name", "description", "price", "active", "created_at", "updated_at")
	productSql = productSql.Values(product.Id, product.Name, product.Description, product.Price, product.Active, product.CreatedAt, product.UpdatedAt)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return entities.Product{}, err
	}

	query, args, err := productSql.ToSql()
	if err != nil {
		tx.Rollback()
		return entities.Product{}, err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		tx.Rollback()
		return entities.Product{}, err
	}

	sql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	newSql := sql.Insert("products_categories").Columns("product_id", "category_id")
	for _, categoryId := range product.CategoriesId {
		newSql = newSql.Values(product.Id, categoryId)
	}
	query, args, err = newSql.ToSql()
	if err != nil {
		tx.Rollback()
		return entities.Product{}, err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		_ = tx.Rollback()
		return entities.Product{}, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return entities.Product{}, err
	}
	return product, nil
}

func (r ProductRepositoryPostgres) UpdateProductFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) (entities.Product, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updateSql := psql.Update("products")
	for key, value := range fields {
		updateSql = updateSql.Set(key, value)
	}
	updateSql = updateSql.Where("id = ?", id)
	query, args, err := updateSql.ToSql()
	if err != nil {
		return entities.Product{}, err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return entities.Product{}, err
	}
	product, err := r.GetProductById(ctx, id)
	if err != nil {
		return entities.Product{}, nil
	}
	return product, nil
}
