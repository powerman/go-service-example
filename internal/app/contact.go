package app

func (a *App) Contacts(ctx Ctx, auth Auth) ([]Contact, error) {
	return a.repo.Contacts(ctx)
}

func (a *App) AddContact(ctx Ctx, auth Auth, name string) (*Contact, error) {
	c := &Contact{
		Name: name,
	}
	err := a.repo.AddContact(ctx, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
