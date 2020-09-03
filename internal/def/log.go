package def

import (
	"github.com/powerman/structlog"
)

// Log field names.
const (
	LogRemote     = "remote" // aligned IPv4:Port "   192.168.0.42:1234 "
	LogFunc       = "func"   // RPC/event handler method name, REST resource path
	LogHost       = "host"
	LogPort       = "port"
	LogAddr       = "addr"
	LogService    = "service"
	LogUser       = "userID"
	LogHTTPMethod = "httpMethod"
	LogHTTPStatus = "httpStatus"
)

func setupLog() {
	structlog.DefaultLogger.
		AppendPrefixKeys(
			LogRemote,
			LogHTTPStatus,
			LogHTTPMethod,
			LogFunc,
		).
		SetSuffixKeys(
			LogService,
			LogUser,
			structlog.KeyStack,
		).
		SetDefaultKeyvals(
			structlog.KeyPID, nil,
		).
		SetKeysFormat(map[string]string{
			structlog.KeyApp:  " %12.12[2]s:", // set to max microservice name length
			structlog.KeyUnit: " %9.9[2]s:",   // set to max KeyUnit/package length
			LogRemote:         " %-21[2]s",
			LogHTTPStatus:     " %3[2]v",
			LogHTTPMethod:     " %-7[2]s",
			LogFunc:           " %[2]s:",
			LogHost:           " %[2]s",
			LogPort:           ":%[2]v",
			LogAddr:           " %[2]s",
			"version":         " %s %v",
			"json":            " %s=%#q",
			"ptr":             " %[2]p",   // for debugging references
			"data":            " %#+[2]v", // for debugging structs
			"offset":          " page=%3[2]d",
			"limit":           "+%[2]d ",
			"err":             " %s: %v",
			LogService:        " [%[2]s]",
			LogUser:           " %[2]v",
		})
}
