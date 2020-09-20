package openapi

import (
	"github.com/powerman/go-service-example/api/openapi/model"
	"github.com/powerman/go-service-example/internal/app"
)

func apiContact(v app.Contact) *model.Contact {
	return &model.Contact{
		ID:   int32(v.ID),
		Name: &v.Name,
	}
}

func apiContacts(vs []app.Contact) []*model.Contact {
	ms := make([]*model.Contact, len(vs))
	for i := range vs {
		ms[i] = apiContact(vs[i])
	}
	return ms
}

func appSeekPage(m model.SeekPagination) app.SeekPage {
	return app.SeekPage{
		SinceID: int(*m.SinceID),
		Limit:   int(*m.Limit),
	}
}
