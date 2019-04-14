// Package def provides default values for both commands and tests.
package def

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/swag"
	"github.com/powerman/go-service-goswagger-clean-example/internal/api/restapi"
	"github.com/powerman/must"
	"github.com/powerman/structlog"
)

// Log field names.
const (
	LogHost   = "host"
	LogPort   = "port"
	LogAddr   = "addr"
	LogRemote = "remote" // aligned IPv4:Port "   192.168.0.42:1234 "
	LogFunc   = "func"   // RPC method name, REST resource path
	LogUser   = "userID"
)

// Default values.
var (
	swaggerHost, swaggerPort = swaggerHostPort()
	Host                     = strGetenv("EXAMPLE_HOST", swaggerHost)
	Port                     = intGetenv("EXAMPLE_PORT", swaggerPort)
	TestTimeFactor           = floatGetenv("GO_TEST_TIME_FACTOR", 1.0)
	TestSecond               = time.Duration(float64(time.Second) * TestTimeFactor)
)

// Init must be called once before using this package.
// It provides common initialization for both commands and tests.
func Init() {
	time.Local = time.UTC
	must.AbortIf = must.PanicIf
	structlog.DefaultLogger.
		AppendPrefixKeys(
			LogRemote,
			LogFunc,
		).
		SetSuffixKeys(
			LogUser,
			structlog.KeyStack,
		).
		SetKeysFormat(map[string]string{
			structlog.KeyUnit: " %6[2]s:", // set to max KeyUnit/package length
			LogHost:           " %[2]s",
			LogPort:           ":%[2]v",
			LogAddr:           " %[2]s",
			LogRemote:         " %-21[2]s",
			LogFunc:           " %[2]s:",
			LogUser:           " %[2]v",
			"version":         " %s %v",
			"err":             " %s: %v",
			"json":            " %s=%#q",
			"ptr":             " %[2]p", // for debugging references
		})
}

func floatGetenv(name string, def float64) float64 {
	value := os.Getenv(name)
	if value == "" {
		return def
	}
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return def
	}
	return v
}

func intGetenv(name string, def int) int {
	value := os.Getenv(name)
	if value == "" {
		return def
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return def
	}
	return i
}

func strGetenv(name, def string) string {
	value := os.Getenv(name)
	if value == "" {
		return def
	}
	return value
}

func swaggerHostPort() (host string, port int) {
	const portHTTP = 80
	const portHTTPS = 443

	spec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return "", 0
	}

	host, port, err = swag.SplitHostPort(spec.Host())
	switch {
	case err == nil:
		return host, port
	case strings.Contains(err.Error(), "missing port"):
		schemes := spec.Spec().Schemes
		switch {
		case len(schemes) == 1 && schemes[0] == "http":
			return spec.Host(), portHTTP
		case len(schemes) == 1 && schemes[0] == "https":
			return spec.Host(), portHTTPS
		}
	}
	return spec.Host(), 0
}
