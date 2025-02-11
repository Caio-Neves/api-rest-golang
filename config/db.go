package config

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"
)

type SqlServerDBConfig struct {
	User     string `toml:"user"`
	Pass     string `toml:"pass"`
	Port     string `toml:"port"`
	Database string `toml:"database"`
	DbServer string `toml:"dbServer"`
}

type PostgresSqlDBConfig struct {
	User     string `toml:"user"`
	Pass     string `toml:"pass"`
	Port     string `toml:"port"`
	Database string `toml:"database"`
	Host     string `toml:"host"`
}

func NewDatabaseConnectionSqlServer(cfg SqlServerDBConfig) (*sql.DB, error) {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s", strings.ReplaceAll(cfg.DbServer, "/", "\\"), cfg.User, cfg.Pass, cfg.Port)
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(5)
	conn.SetMaxIdleConns(1)
	return conn, nil
}

func NewDatabaseConnectionPostgreSQL(cfg PostgresSqlDBConfig) (*sql.DB, error) {
	connString := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable dbname=%s", cfg.Host, cfg.User, cfg.Pass, cfg.Port, cfg.Database)
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(5)
	conn.SetMaxIdleConns(1)
	return conn, nil

}
