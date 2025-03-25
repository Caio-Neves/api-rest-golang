package user

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"rest-api-example/entities"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrErrorGenerateHash  = errors.New("could not generate password hash")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPassword    = errors.New("invalid password")
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) entities.UserInterface {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) CheckUserCredentials(ctx context.Context, credentials entities.Credentials) (bool, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	categorySql := psql.Select("login", "password").From("users").Where(sq.Eq{"login": credentials.Login})
	query, args, err := categorySql.ToSql()
	if err != nil {
		return false, err
	}
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	c := entities.Credentials{}
	row := stmt.QueryRowContext(ctx, args...)
	if errors.Is(row.Err(), sql.ErrNoRows) {
		return false, ErrUserNotFound
	}
	err = row.Scan(&c.Login, &c.Password)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(credentials.Password))
	if err != nil {
		return false, ErrInvalidPassword
	}
	return true, nil
}

func (r UserRepository) RegistryUser(ctx context.Context, credentials entities.Credentials) error {
	pass, _, err := r.hashUserPassword(credentials)
	if err != nil {
		return err
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	usersSql := psql.Insert("users").Columns("login", "password")
	usersSql = usersSql.Values(credentials.Login, pass)

	query, args, err := usersSql.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

type hash string
type salt string

func (r UserRepository) saltUserPassword(credentials entities.Credentials) salt {
	return salt(fmt.Sprintf("%s.%s*%s", "9aca55433a5217bb", credentials.Login, credentials.Password))
}

func (r UserRepository) hashUserPassword(credentials entities.Credentials) (hash, salt, error) {
	salted := r.saltUserPassword(credentials)
	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(salted), bcrypt.DefaultCost)
	if err != nil {
		return "", "", ErrErrorGenerateHash
	}
	return hash(hex.EncodeToString(bcryptHash)), salted, nil
}
