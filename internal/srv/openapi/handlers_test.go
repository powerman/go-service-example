package openapi_test

import (
	"io"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/powerman/go-service-goswagger-clean-example/api/openapi/client/op"
	"github.com/powerman/go-service-goswagger-clean-example/api/openapi/model"
	"github.com/powerman/go-service-goswagger-clean-example/internal/app"
	"github.com/powerman/go-service-goswagger-clean-example/internal/srv/openapi"
)

func TestListContacts(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	c, _, mockApp := testNewServer(t)
	params := op.NewListContactsParams()

	mockApp.EXPECT().Contacts(gomock.Any(), authUser).Return(nil, io.EOF)
	mockApp.EXPECT().Contacts(gomock.Any(), authUser).Return(nil, nil)
	mockApp.EXPECT().Contacts(gomock.Any(), authAdmin).Return([]app.Contact{appContact1, appContact2}, nil)

	testCases := []struct {
		apiKey  runtime.ClientAuthInfoWriter
		want    []*model.Contact
		wantErr *model.Error
	}{
		{nil, nil, apiError401},
		{apiKeyUser, nil, apiError500},
		{apiKeyUser, []*model.Contact{}, nil},
		{apiKeyAdmin, []*model.Contact{apiContact1, apiContact2}, nil},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := c.Op.ListContacts(params, tc.apiKey)
			if tc.wantErr == nil {
				t.Nil(openapi.ErrPayload(err))
				t.DeepEqual(res.Payload, tc.want)
			} else {
				t.DeepEqual(openapi.ErrPayload(err), tc.wantErr)
				t.Nil(res)
			}
		})
	}
}

func TestAddContact(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	c, _, mockApp := testNewServer(t)
	params := op.NewAddContactParams()

	mockApp.EXPECT().AddContact(gomock.Any(), authAdmin, " ").Return(nil, io.EOF)
	mockApp.EXPECT().AddContact(gomock.Any(), authAdmin, "A").Return(nil, app.ErrContactExists)
	mockApp.EXPECT().AddContact(gomock.Any(), authAdmin, "B").Return(&appContact2, nil)

	testCases := []struct {
		apiKey  runtime.ClientAuthInfoWriter
		name    string
		want    *model.Contact
		wantErr *model.Error
	}{
		{nil, "A", nil, apiError401},
		{apiKeyUser, "A", nil, apiError403},
		{apiKeyAdmin, " ", nil, apiError500},
		{apiKeyAdmin, "A", nil, apiError1000},
		{apiKeyAdmin, "B", apiContact2, nil},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			params.Contact = &model.Contact{Name: swag.String(tc.name)}
			res, err := c.Op.AddContact(params, tc.apiKey)
			if tc.wantErr == nil {
				t.Nil(openapi.ErrPayload(err))
				t.DeepEqual(res.Payload, tc.want)
			} else {
				t.DeepEqual(openapi.ErrPayload(err), tc.wantErr)
				t.Nil(res)
			}
		})
	}
}
