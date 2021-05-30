package app_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/powerman/check"

	"github.com/powerman/go-service-example/internal/app"
)

func TestContacts(tt *testing.T) {
	t := check.T(tt)
	a, mockRepo := testNew(t)

	var (
		c1    = app.Contact{ID: 1, Name: "A"}
		c2    = app.Contact{ID: 2, Name: "B"}
		page1 = app.SeekPage{SinceID: 0, Limit: 2}
		page2 = app.SeekPage{SinceID: 2, Limit: 2}
	)

	mockRepo.EXPECT().LstContacts(gomock.Any(), page1).Return([]app.Contact{c1, c2}, nil)
	mockRepo.EXPECT().LstContacts(gomock.Any(), page2).Return(nil, nil)

	db, err := a.Contacts(ctx, auth1, page1)
	t.Nil(err)
	t.Len(db, 2)

	db, err = a.Contacts(ctx, auth1, page2)
	t.Nil(err)
	t.Zero(db)
}

func TestAddContact(tt *testing.T) {
	t := check.T(tt)
	a, mockRepo := testNew(t)

	c1 := app.Contact{Name: "A"}

	mockRepo.EXPECT().AddContact(gomock.Any(), c1.Name).Return(1, nil)
	mockRepo.EXPECT().AddContact(gomock.Any(), c1.Name).Return(0, app.ErrContactExists)

	c, err := a.AddContact(ctx, auth1, c1.Name)
	t.Nil(err)
	t.DeepEqual(c, &app.Contact{ID: 1, Name: c1.Name})

	c, err = a.AddContact(ctx, auth1, c1.Name)
	t.Err(err, app.ErrContactExists)
	t.Nil(c)
}
