//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mock.$GOFILE Appl,Repo

// Package app provides business logic.
package app

import (
	"context"
	"errors"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

// Errors.
var (
	ErrContactExists = errors.New("contact already exists")
)

// Appl provides application features (use cases) service.
type Appl interface {
	// HealthCheck returns error if service is unhealthy or current
	// status otherwise.
	// Errors: none.
	HealthCheck(Ctx) (interface{}, error)
	// Contacts returns all contacts.
	// Errors: none.
	Contacts(Ctx, Auth, SeekPage) ([]Contact, error)
	// AddContact adds new contact.
	// Errors: ErrContactExists.
	AddContact(_ Ctx, _ Auth, name string) (*Contact, error)
}

// Repo provides data storage.
type Repo interface {
	// LstContacts returns up to limit contacts with ID > sinceID,
	// ordered by ID.
	// Errors: none.
	LstContacts(Ctx, SeekPage) ([]Contact, error)
	// AddContact adds new contact and returns it ID.
	// Errors: ErrContactExists.
	AddContact(_ Ctx, name string) (id int, err error)
}

type (
	// Auth describes authentication.
	Auth struct {
		UserID string
	}
	// SeekPage describes seek pagination.
	SeekPage struct {
		SinceID int
		Limit   int
	}
	// Contact describes record in address book.
	Contact struct {
		ID   int
		Name string
	}
)

// App implements interface Appl.
type App struct {
	repo Repo
}

func New(repo Repo) *App {
	a := &App{
		repo: repo,
	}
	return a
}

func (a *App) HealthCheck(_ Ctx) (interface{}, error) {
	return "OK", nil
}
