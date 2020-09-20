package dal_test

import (
	"testing"

	"github.com/powerman/check"
	"github.com/powerman/go-service-example/internal/app"
	dal "github.com/powerman/go-service-example/internal/dal/memory"
)

func TestContacts(tt *testing.T) {
	t := check.T(tt)
	r, err := dal.New(ctx)
	t.Nil(err)

	db, err := r.Contacts(ctx)
	t.Nil(err)
	t.Zero(db)

	c := &app.Contact{Name: "A"}
	err = r.AddContact(ctx, c)
	t.Nil(err)
	t.Equal(c.ID, 1)

	c = &app.Contact{Name: "A"}
	err = r.AddContact(ctx, c)
	t.Err(err, app.ErrContactExists)
	t.Zero(c.ID)

	c = &app.Contact{Name: "B"}
	err = r.AddContact(ctx, c)
	t.Nil(err)
	t.Equal(c.ID, 2)

	db, err = r.Contacts(ctx)
	t.Nil(err)
	t.Len(db, 2)
}
