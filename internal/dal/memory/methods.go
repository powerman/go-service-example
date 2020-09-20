package dal

import (
	"github.com/powerman/structlog"

	"github.com/powerman/go-service-example/internal/app"
)

func (r *Repo) AddContact(ctx Ctx, name string) (id int, err error) {
	log := structlog.FromContext(ctx, nil)
	r.Lock()
	defer r.Unlock()

	r.lastID++

	for i := range r.db {
		if r.db[i].Name == name {
			return 0, app.ErrContactExists
		}
	}

	id = r.lastID
	r.db = append(r.db, app.Contact{ID: id, Name: name})
	log.Debug("contact added")
	return id, nil
}

func (r *Repo) LstContacts(ctx Ctx, page app.SeekPage) ([]app.Contact, error) {
	r.Lock()
	defer r.Unlock()

	contacts := make([]app.Contact, 0, page.Limit)
	for i := range r.db {
		if len(contacts) >= page.Limit {
			break
		}
		if r.db[i].ID > page.SinceID {
			contacts = append(contacts, r.db[i])
		}
	}
	return contacts, nil
}
