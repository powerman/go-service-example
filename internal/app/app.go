// Package app provides business logic.
package app

import (
	"context"

	"github.com/powerman/structlog"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

// Log is a synonym for convenience.
type Log = *structlog.Logger

// Auth describes authentication.
type Auth struct {
	UserID string
}

// App provides application features service.
type App interface {
	Contacts(Ctx, Log, Auth) ([]Contact, error)
	AddContact(Ctx, Log, Auth, *Contact) error
}

type app struct {
	lastID int
	db     []Contact
}

// New return new application.
func New() App {
	return &app{}
}
