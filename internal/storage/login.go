package storage

import (
	"context"

	mssql "github.com/microsoft/go-mssqldb"
	"github.com/pkg/errors"
)

func (s *Service) GetUserPassword(ctx context.Context, login string) (password string, err error) {
	sql := "SELECT password FROM Employee WHERE Login = @p1"

	err = s.DB.QueryRowContext(ctx, sql, mssql.VarChar(login)).Scan(&password)
	if err != nil {
		return "", errors.Wrap(err, "failed to query for user")
	}

	return password, nil
}
