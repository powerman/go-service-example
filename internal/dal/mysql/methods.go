package dal

import (
	"github.com/powerman/go-service-example/internal/app"
	"github.com/powerman/go-service-example/pkg/repo"
)

func (r *Repo) AddContact(ctx Ctx, name string) (id int, err error) {
	err = r.NoTx(func() error {
		res, err := r.DB.NamedExecContext(ctx, sqlContactAdd, argContactAdd{
			Name: name,
		})
		switch {
		case repo.MySQLDuplicateEntry(err):
			return app.ErrContactExists
		case err != nil:
			return err
		default:
			insertID, err := res.LastInsertId()
			id = int(insertID)
			return err
		}
	})
	return
}

func (r *Repo) LstContacts(ctx Ctx, page app.SeekPage) (contacts []app.Contact, err error) {
	err = r.NoTx(func() error {
		var rows []rowContactLst
		err := r.DB.NamedSelectContext(ctx, &rows, sqlContactLst, argContactLst{
			SinceID: page.SinceID,
			Limit:   page.Limit,
		})
		if err != nil {
			return err
		}
		contacts = make([]app.Contact, len(rows))
		for i := range rows {
			contacts[i] = appContact(rows[i])
		}
		return nil
	})
	return
}
