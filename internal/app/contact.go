package app

func (a *App) Contacts(ctx Ctx, auth Auth, page SeekPage) ([]Contact, error) {
	return a.repo.LstContacts(ctx, page)
}

func (a *App) AddContact(ctx Ctx, auth Auth, name string) (*Contact, error) {
	id, err := a.repo.AddContact(ctx, name)
	if err != nil {
		return nil, err
	}
	c := &Contact{
		ID:   id,
		Name: name,
	}
	return c, nil
}
