package def

import (
	"github.com/powerman/structlog"
)

// Log field names.
const (
	LogServer     = "server"     // "OpenAPI", "gRPC", "Prometheus metrics", etc.
	LogRemoteIP   = "remoteIP"   // IP address.
	LogAddr       = "addr"       // host:port.
	LogHost       = "host"       // DNS hostname or IPv4/IPv6 address.
	LogPort       = "port"       // TCP/UDP port number.
	LogHTTPMethod = "httpMethod" // GET, POST, etc.
	LogHTTPStatus = "httpStatus" // Status code: 200, 404, etc.
	LogFunc       = "func"       // RPC/event handler method name, REST resource path.
	LogUserID     = "userID"
)

func setupLog() {
	structlog.DefaultLogger.
		AppendPrefixKeys(
			LogRemoteIP,
			LogHTTPStatus,
			LogHTTPMethod,
			LogFunc,
		).
		SetSuffixKeys(
			LogServer,
			LogUserID,
			"dump",
			structlog.KeyStack,
		).
		SetDefaultKeyvals(
			structlog.KeyPID, nil,
		).
		SetKeysFormat(map[string]string{
			structlog.KeyApp:  " %12.12[2]s:", // set to max microservice name length
			structlog.KeyUnit: " %9.9[2]s:",   // set to max KeyUnit/package length
			LogRemoteIP:       " %-15[2]s",    // set to 19.19 or 39 or 45 for IPv6
			LogHTTPStatus:     " %3[2]v",
			LogHTTPMethod:     " %-7[2]s",
			LogFunc:           " %[2]s:",
			LogHost:           " %[2]s",
			LogPort:           ":%[2]v",
			LogAddr:           " %[2]s",
			"version":         " %s %v",
			"json":            " %s=%#q",
			"ptr":             " %[2]p",            // for debugging references
			"data":            " %#+[2]v",          // for debugging structs
			"dump":            "\n›››\n%[2]s\n‹‹‹", // for debugging multiline text
			"offset":          " page=%3[2]d",
			"limit":           "+%[2]d ",
			"err":             " %s: %v",
			LogServer:         " [%[2]s]",
			LogUserID:         " %[2]v",
		})
}
