package openapi

import (
	"github.com/powerman/go-service-goswagger-clean-example/api/openapi/model"
	"github.com/powerman/go-service-goswagger-clean-example/internal/app"
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
