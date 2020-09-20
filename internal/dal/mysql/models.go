package dal

import "github.com/powerman/go-service-example/internal/app"

func appContact(v rowContactLst) app.Contact {
	return app.Contact{
		ID:   v.ID,
		Name: v.Name,
	}
}
