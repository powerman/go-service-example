package openapi_test

import (
	"io"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/golang/mock/gomock"
	"github.com/powerman/check"

	"github.com/powerman/go-service-example/api/openapi/client/op"
	"github.com/powerman/go-service-example/api/openapi/model"
	"github.com/powerman/go-service-example/internal/app"
	"github.com/powerman/go-service-example/internal/srv/openapi"
)

func TestHealthCheck(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	c, _, mockAppl, _ := testNewServer(t, openapi.Config{})

	mockAppl.EXPECT().HealthCheck(gomock.Any()).Return(nil, io.EOF)
	mockAppl.EXPECT().HealthCheck(gomock.Any()).Return(nil, nil)
	mockAppl.EXPECT().HealthCheck(gomock.Any()).Return("OK", nil)
	mockAppl.EXPECT().HealthCheck(gomock.Any()).Return(map[string]string{"main": "OK"}, nil)

	testCases := []struct {
		want    interface{}
		wantErr *model.Error
	}{
		{nil, apiError500},
		{nil, nil},
		{"OK", nil},
		{map[string]interface{}{"main": "OK"}, nil},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := c.Op.HealthCheck(op.NewHealthCheckParams())
			t.DeepEqual(openapi.ErrPayload(err), tc.wantErr)
			if res == nil {
				t.DeepEqual(nil, tc.want)
			} else {
				t.DeepEqual(res.Payload, tc.want)
			}
		})
	}
}

func TestListContacts(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	c, _, mockAppl, _ := testNewServer(t, openapi.Config{})

	var (
		apiPage1 = model.SeekPagination{SinceID: swag.Int32(0), Limit: swag.Int32(2)}
		apiPage2 = model.SeekPagination{SinceID: swag.Int32(2), Limit: swag.Int32(2)}
		appPage1 = app.SeekPage{SinceID: 0, Limit: 2}
		appPage2 = app.SeekPage{SinceID: 2, Limit: 2}
	)

	mockAppl.EXPECT().Contacts(gomock.Any(), authUser, appPage1).Return(nil, io.EOF)
	mockAppl.EXPECT().Contacts(gomock.Any(), authUser, appPage2).Return(nil, nil)
	mockAppl.EXPECT().Contacts(gomock.Any(), authAdmin, appPage1).Return([]app.Contact{appContact1, appContact2}, nil)

	testCases := []struct {
		apiKey  runtime.ClientAuthInfoWriter
		page    model.SeekPagination
		want    interface{}
		wantErr *model.Error
	}{
		{nil, apiPage1, nil, apiError401},
		{apiKeyUser, apiPage1, nil, apiError500},
		{apiKeyUser, apiPage2, []*model.Contact{}, nil},
		{apiKeyAdmin, apiPage1, []*model.Contact{apiContact1, apiContact2}, nil},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			args := op.ListContactsBody{SeekPagination: tc.page}
			res, err := c.Op.ListContacts(op.NewListContactsParams().WithArgs(args), tc.apiKey)
			t.DeepEqual(openapi.ErrPayload(err), tc.wantErr)
			if res == nil {
				t.DeepEqual(nil, tc.want)
			} else {
				t.DeepEqual(res.Payload, tc.want)
			}
		})
	}
}

func TestAddContact(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	c, _, mockAppl, _ := testNewServer(t, openapi.Config{})

	mockAppl.EXPECT().AddContact(gomock.Any(), authAdmin, " ").Return(nil, io.EOF)
	mockAppl.EXPECT().AddContact(gomock.Any(), authAdmin, "A").Return(nil, app.ErrContactExists)
	mockAppl.EXPECT().AddContact(gomock.Any(), authAdmin, "B").Return(&appContact2, nil)

	testCases := []struct {
		apiKey  runtime.ClientAuthInfoWriter
		name    string
		want    interface{}
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
			args := &model.Contact{Name: swag.String(tc.name)}
			res, err := c.Op.AddContact(op.NewAddContactParams().WithArgs(args), tc.apiKey)
			t.DeepEqual(openapi.ErrPayload(err), tc.wantErr)
			if res == nil {
				t.DeepEqual(nil, tc.want)
			} else {
				t.DeepEqual(res.Payload, tc.want)
			}
		})
	}
}
