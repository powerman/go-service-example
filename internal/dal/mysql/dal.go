// Package dal implements Data Access Layer using MySQL DB.
package dal

import (
	"context"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/powerman/go-service-example/internal/app"
	migrations "github.com/powerman/go-service-example/internal/migrations/mysql"
	"github.com/powerman/go-service-example/pkg/repo"
)

const (
	schemaVersion  = 3
	dbMaxOpenConns = 0 // Unlimited.
	dbMaxIdleConns = 5 // A bit more than default (2).
)

type Ctx = context.Context

// Repo provides access to storage.
type Repo struct {
	*repo.Repo
}

// New creates and returns new Repo.
// It will also run required DB migrations and connects to DB.
func New(ctx Ctx, dir string, cfg *mysql.Config) (_ *Repo, err error) {
	returnErrs := []error{ // List of app.Err… returned by Repo methods.
		app.ErrContactExists,
	}

	r := &Repo{}
	r.Repo, err = repo.New(ctx, migrations.Goose(), repo.Config{
		MySQL:         cfg,
		GooseDir:      dir,
		SchemaVersion: schemaVersion,
		Metric:        metric,
		ReturnErrs:    returnErrs,
	})
	if err != nil {
		return nil, err
	}
	r.DB.SetMaxOpenConns(dbMaxOpenConns)
	r.DB.SetMaxIdleConns(dbMaxIdleConns)
	r.SchemaVer.HoldSharedLock(ctx, time.Second)
	return r, nil
}
