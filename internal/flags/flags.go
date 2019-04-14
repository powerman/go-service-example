// Package flags provide helpers to use with package flag.
package flags

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
)

// FatalFlagValue report invalid flag values in same way as flag.Parse().
func FatalFlagValue(msg, name string, val interface{}) {
	fmt.Fprintf(os.Stderr, "invalid value %#v for flag -%s: %s\n", val, name, msg)
	flag.Usage()
	os.Exit(2)
}

// Endpoint validates url and ensure it doesn't have trailing slashes.
func Endpoint(s *string) bool {
	clean := regexp.MustCompile(`/+$`).ReplaceAllString(*s, "")

	p, err := url.Parse(clean)
	if err != nil {
		return false
	}
	if p.Host == "" {
		return false
	}

	*s = clean
	return true
}
