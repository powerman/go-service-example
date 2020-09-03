package app_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/powerman/go-service-goswagger-clean-example/internal/app"
)

func TestContacts(tt *testing.T) {
	t := check.T(tt)
	a, mockRepo := testNew(t)

	var (
		c1 = app.Contact{ID: 1, Name: "A"}
		c2 = app.Contact{ID: 2, Name: "B"}
	)

	mockRepo.EXPECT().Contacts(gomock.Any()).Return(nil, nil)
	mockRepo.EXPECT().Contacts(gomock.Any()).Return([]app.Contact{c1, c2}, nil)

	db, err := a.Contacts(ctx, auth1)
	t.Nil(err)
	t.Zero(db)

	db, err = a.Contacts(ctx, auth1)
	t.Nil(err)
	t.Len(db, 2)
}

func TestAddContact(tt *testing.T) {
	t := check.T(tt)
	a, mockRepo := testNew(t)

	c1 := app.Contact{Name: "A"}

	mockRepo.EXPECT().AddContact(gomock.Any(), &c1).DoAndReturn(func(_ Ctx, c *app.Contact) error {
		c.ID = 1
		return nil
	})
	mockRepo.EXPECT().AddContact(gomock.Any(), &c1).Return(app.ErrContactExists)

	c, err := a.AddContact(ctx, auth1, c1.Name)
	t.Nil(err)
	t.DeepEqual(c, &app.Contact{ID: 1, Name: c1.Name})

	c, err = a.AddContact(ctx, auth1, c1.Name)
	t.Err(err, app.ErrContactExists)
	t.Nil(c)
}
