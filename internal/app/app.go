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
	// Contacts returns all contacts.
	// Errors: none.
	Contacts(Ctx, Auth) ([]Contact, error)
	// AddContact adds new contact.
	// Errors: ErrContactExists.
	AddContact(_ Ctx, _ Auth, name string) (*Contact, error)
}

// Repo provides data storage.
type Repo interface {
	// Contacts returns all contacts.
	// Errors: none.
	Contacts(Ctx) ([]Contact, error)
	// AddContact adds new contact and set ID.
	// Errors: ErrContactExists.
	AddContact(Ctx, *Contact) error
}

type (
	// Auth describes authentication.
	Auth struct {
		UserID string
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
