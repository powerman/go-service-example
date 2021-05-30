// +build integration

package dal_test

import (
	"testing"

	"github.com/powerman/check"

	"github.com/powerman/go-service-example/internal/app"
)

func TestContact(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	r := newTestRepo(t)

	var (
		c1 = app.Contact{ID: 1, Name: "A"}
		c3 = app.Contact{ID: 3, Name: "B"}
		c4 = app.Contact{ID: 4, Name: "C"}
	)

	contacts, err := r.LstContacts(ctx, app.SeekPage{SinceID: 0, Limit: 2})
	t.Nil(err)
	t.Len(contacts, 0)

	testsAdd := []struct {
		name    string
		want    int
		wantErr error
	}{
		{c1.Name, c1.ID, nil},
		{c1.Name, 0, app.ErrContactExists},
		{c3.Name, c3.ID, nil},
		{c4.Name, c4.ID, nil},
	}
	for _, tc := range testsAdd {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := r.AddContact(ctx, tc.name)
			t.Err(err, tc.wantErr)
			t.Equal(res, tc.want)
		})
	}

	testsLst := []struct {
		page    app.SeekPage
		want    []app.Contact
		wantErr error
	}{
		{app.SeekPage{SinceID: 0, Limit: 0}, []app.Contact{}, nil},
		{app.SeekPage{SinceID: 0, Limit: 2}, []app.Contact{c1, c3}, nil},
		{app.SeekPage{SinceID: 2, Limit: 5}, []app.Contact{c3, c4}, nil},
		{app.SeekPage{SinceID: c3.ID, Limit: 2}, []app.Contact{c4}, nil},
		{app.SeekPage{SinceID: c4.ID, Limit: 2}, []app.Contact{}, nil},
	}
	for _, tc := range testsLst {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := r.LstContacts(ctx, tc.page)
			t.Err(err, tc.wantErr)
			t.DeepEqual(res, tc.want)
		})
	}
}
