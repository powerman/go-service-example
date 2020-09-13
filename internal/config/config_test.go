package config

import (
	"os"
	"testing"

	"github.com/powerman/check"
	"github.com/powerman/go-service-example/pkg/def"
	"github.com/powerman/go-service-example/pkg/netx"
)

func Test(t *testing.T) {
	want := &ServeConfig{
		APIKeyAdmin: "admin",
		Addr:        netx.NewAddr(def.Hostname, 8000),
		MetricsAddr: netx.NewAddr(def.Hostname, 9000),
	}

	t.Run("required", func(tt *testing.T) {
		t := check.T(tt)
		require(t, "APIKeyAdmin")
		os.Setenv("EXAMPLE_APIKEY_ADMIN", "admin")
	})
	t.Run("default", func(tt *testing.T) {
		t := check.T(tt)
		c, err := testGetServe()
		t.Nil(err)
		t.DeepEqual(c, want)
	})
	t.Run("constraint", func(tt *testing.T) {
		t := check.T(tt)
		constraint(t, "EXAMPLE_ADDR_PORT", "x", `^AddrPort .* invalid syntax`)
		constraint(t, "EXAMPLE_METRICS_ADDR_PORT", "x", `^MetricsAddrPort .* invalid syntax`)
	})
	t.Run("env", func(tt *testing.T) {
		t := check.T(tt)
		os.Setenv("EXAMPLE_APIKEY_ADMIN", "admin3")
		os.Setenv("EXAMPLE_ADDR_HOST", "localhost3")
		os.Setenv("EXAMPLE_ADDR_PORT", "8003")
		os.Setenv("EXAMPLE_METRICS_ADDR_PORT", "9003")
		c, err := testGetServe()
		t.Nil(err)
		want.APIKeyAdmin = "admin3"
		want.Addr = netx.NewAddr("localhost3", 8003)
		want.MetricsAddr = netx.NewAddr("localhost3", 9003)
		t.DeepEqual(c, want)
	})
	t.Run("flag", func(tt *testing.T) {
		t := check.T(tt)
		c, err := testGetServe(
			"--host=localhost4",
			"--port=8004",
			"--metrics.port=9004",
		)
		t.Nil(err)
		want.Addr = netx.NewAddr("localhost4", 8004)
		want.MetricsAddr = netx.NewAddr("localhost4", 9004)
		t.DeepEqual(c, want)
	})
	t.Run("cleanup", func(tt *testing.T) {
		t := check.T(tt)
		t.Panic(func() { GetServe() })
	})
}
