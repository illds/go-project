package db

import (
	"GOHW-1/internal/configuration"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewDb(ctx context.Context, dbCredentials *configuration.DBCredentials) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, generateDsn(dbCredentials))
	if err != nil {
		return nil, err
	}
	return newDataBase(pool), nil
}

func generateDsn(dbCredentials *configuration.DBCredentials) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbCredentials.Host, dbCredentials.Port, dbCredentials.User, dbCredentials.Password, dbCredentials.DBname)
}
