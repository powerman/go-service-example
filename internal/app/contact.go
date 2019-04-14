package app

// Contact describes record in address book.
type Contact struct {
	ID   int
	Name string
}

func (app *app) Contacts(ctx Ctx, log Log, auth Auth) ([]Contact, error) {
	return app.db, nil
}

func (app *app) AddContact(ctx Ctx, log Log, auth Auth, c *Contact) error {
	app.lastID++
	c.ID = app.lastID
	app.db = append(app.db, *c)
	return nil
}
