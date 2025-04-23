package user

import (
	"context"
	"database/sql"
	"errors"
	"rest-api-example/entities"

	sq "github.com/Masterminds/squirrel"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) entities.UserInterface {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) GetCredentialsByLogin(ctx context.Context, login string) (entities.Credentials, error) {
	op := "UserRepository.GetCredentials()"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	categorySql := psql.Select("login", "password").From("users").Where(sq.Eq{"login": login})
	query, args, err := categorySql.ToSql()
	if err != nil {
		return entities.Credentials{}, entities.NewInternalServerErrorError(err, op)
	}

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return entities.Credentials{}, entities.NewInternalServerErrorError(err, op)
	}
	defer stmt.Close()

	credentials := entities.Credentials{}
	row := stmt.QueryRowContext(ctx, args...)
	if errors.Is(row.Err(), sql.ErrNoRows) {
		return entities.Credentials{}, nil
	}
	err = row.Scan(&credentials.Login, &credentials.Password)
	if err != nil {
		return entities.Credentials{}, entities.NewInternalServerErrorError(err, op)
	}
	return credentials, nil
}

func (r UserRepository) InsertUser(ctx context.Context, credentials entities.Credentials) error {
	op := "UserRepository.InsertUser()"
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	usersSql := psql.Insert("users").Columns("login", "password")
	usersSql = usersSql.Values(credentials.Login, credentials.Password)

	query, args, err := usersSql.ToSql()
	if err != nil {
		return entities.NewInternalServerErrorError(err, op)
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return entities.NewInternalServerErrorError(err, op)
	}
	return nil
}
