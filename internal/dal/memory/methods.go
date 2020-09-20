package dal

import (
	"github.com/powerman/go-service-example/internal/app"
	"github.com/powerman/structlog"
)

func (r *Repo) Contacts(_ Ctx) ([]app.Contact, error) {
	r.Lock()
	defer r.Unlock() //nolint:gocritic // False positive (unnecessaryDefer).

	return r.db, nil
}

func (r *Repo) AddContact(ctx Ctx, c *app.Contact) error {
	log := structlog.FromContext(ctx, nil)
	r.Lock()
	defer r.Unlock()

	for i := range r.db {
		if r.db[i].Name == c.Name {
			return app.ErrContactExists
		}
	}

	r.lastID++
	c.ID = r.lastID
	r.db = append(r.db, *c)
	log.Debug("contact added")
	return nil
}
