package def

import (
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"

	"github.com/powerman/getenv"
)

// Default ports.
const (
	Port        = 8000
	MetricsPort = 9000
)

// Constants.
var (
	ver                   string // Set by ./build script.
	ProgName              = strings.TrimSuffix(path.Base(os.Args[0]), ".test")
	Hostname, hostnameErr = os.Hostname()
	testTimeFactor        = getenv.Float("GO_TEST_TIME_FACTOR", 1.0)
	TestSecond            = time.Duration(float64(time.Second) * testTimeFactor)
)

// Version returns application version based on build info.
func Version() string {
	if bi, ok := debug.ReadBuildInfo(); ok {
		if bi.Main.Version == "(devel)" && ver != "" {
			return ver
		}
		return bi.Main.Version
	}
	return "(test)"
}
