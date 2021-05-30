// +build integration

package main

import (
	"context"
	"testing"

	"github.com/go-openapi/swag"
	"github.com/powerman/check"
	"github.com/powerman/mysqlx"

	"github.com/powerman/go-service-example/api/openapi/client"
	"github.com/powerman/go-service-example/api/openapi/client/op"
	"github.com/powerman/go-service-example/api/openapi/model"
	"github.com/powerman/go-service-example/internal/srv/openapi"
	"github.com/powerman/go-service-example/pkg/def"
	"github.com/powerman/go-service-example/pkg/netx"
)

func TestSmoke(tt *testing.T) {
	t := check.T(tt)

	s := &Service{cfg: cfg}

	tempDBCfg, cleanup, err := mysqlx.EnsureTempDB(tLogger(*t), "", cfg.MySQL)
	cfg.MySQL = tempDBCfg // Assign to cfg and not s.cfg as a reminder: they are the same.
	t.Must(t.Nil(err))
	defer cleanup()

	ctxStartup, cancel := context.WithTimeout(ctx, def.TestTimeout)
	defer cancel()
	ctxShutdown, shutdown := context.WithCancel(ctx)
	errc := make(chan error)
	go func() { errc <- s.RunServe(ctxStartup, ctxShutdown, shutdown) }()
	defer func() {
		shutdown()
		t.Nil(<-errc, "RunServe")
		if s.repo != nil {
			s.repo.Close()
		}
	}()
	t.Must(t.Nil(netx.WaitTCPPort(ctxStartup, cfg.BindAddr), "connect to service"))

	openapiClient := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Schemes:  []string{"http"},
		Host:     cfg.BindAddr.String(),
		BasePath: client.DefaultBasePath,
	})

	{
		args := &model.Contact{Name: apiContact1.Name}
		res, err := openapiClient.Op.AddContact(op.NewAddContactParams().WithArgs(args), apiKeyUser)
		t.DeepEqual(openapi.ErrPayload(err), apiError403)
		t.Nil(res)
	}
	{
		args := &model.Contact{Name: apiContact1.Name}
		res, err := openapiClient.Op.AddContact(op.NewAddContactParams().WithArgs(args), apiKeyAdmin)
		t.Nil(openapi.ErrPayload(err))
		t.DeepEqual(res, &op.AddContactCreated{Payload: apiContact1})
	}
	{
		args := op.ListContactsBody{SeekPagination: model.SeekPagination{
			SinceID: swag.Int32(0),
			Limit:   swag.Int32(2),
		}}
		res, err := openapiClient.Op.ListContacts(op.NewListContactsParams().WithArgs(args), apiKeyAdmin)
		t.Nil(openapi.ErrPayload(err))
		t.DeepEqual(res, &op.ListContactsOK{Payload: []*model.Contact{apiContact1}})
	}
}
