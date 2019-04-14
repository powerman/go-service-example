package api

import (
	"github.com/powerman/go-service-goswagger-clean-example/internal/api/model"
	"github.com/powerman/go-service-goswagger-clean-example/internal/app"
)

func appContact(m *model.Contact) app.Contact {
	return app.Contact{
		ID:   int(m.ID),
		Name: *m.Name,
	}
}

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
