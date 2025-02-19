package repositories

import (
	"context"
	"database/sql"
	"log"
	"rest-api-example/entities"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CategoryRepositoryPostgres struct {
	db *sql.DB
}

func NewCategoryRepositoryPostgres(db *sql.DB) *CategoryRepositoryPostgres {
	return &CategoryRepositoryPostgres{
		db: db,
	}
}

func (r CategoryRepositoryPostgres) GetAllCategories(ctx context.Context, params map[string][]string) ([]entities.Category, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	categoriesSql := psql.Select("id", "name", "description", "active", "created_at", "updated_at").From("categories")
	if value, exists := params["active"]; exists {
		isActive, err := strconv.Atoi(value[0])
		if err != nil {
			log.Println(err)
			return nil, err
		}
		categoriesSql = categoriesSql.Where("active = ?", isActive)
	}
	page := 1
	limit := 10
	if value, exists := params["page"]; exists {
		page, _ = strconv.Atoi(value[0])
	}
	if value, exists := params["limit"]; exists {
		limit, _ = strconv.Atoi(value[0])
	}
	offset := (page - 1) * limit
	categoriesSql = categoriesSql.Limit(uint64(limit)).Offset(uint64(offset))

	query, args, err := categoriesSql.ToSql()
	log.Println(args)
	log.Println(query)
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

func (r CategoryRepositoryPostgres) GetCategoryById(ctx context.Context, id uuid.UUID) (entities.Category, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	categorySql := psql.Select("id", "name", "description", "active", "created_at", "updated_at").From("categories")
	categorySql = categorySql.Where("id = ?", id)

	query, args, err := categorySql.ToSql()
	log.Println(query)
	log.Println(args)
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

func convertUUIDsToStrings(uuids []uuid.UUID) []string {
	strs := make([]string, len(uuids))
	for i, id := range uuids {
		strs[i] = id.String()
	}
	return strs
}

func (r CategoryRepositoryPostgres) GetCategoriesByIds(ctx context.Context, ids []uuid.UUID) ([]entities.Category, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	categorySql := psql.Select("id", "name", "description", "active", "created_at", "updated_at").From("categories")
	categorySql = categorySql.Where(sq.Eq{"id": convertUUIDsToStrings(ids)})

	query, args, err := categorySql.ToSql()
	log.Println(query)
	log.Println(args)
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
