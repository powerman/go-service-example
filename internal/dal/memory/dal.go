// Package dal implements Data Access Layer using in-memory DB.
package dal

import (
	"context"
	"sync"

	"github.com/powerman/go-service-example/internal/app"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

// Repo provides access to storage.
type Repo struct {
	sync.Mutex
	lastID int
	db     []app.Contact
}

// New creates and returns new Repo.
func New(_ Ctx) (*Repo, error) {
	return &Repo{}, nil
}
