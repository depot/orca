//go:build windows
// +build windows

package hosts

import (
	"os"
)

// DefaultPath is the default path to the hosts file.
var DefaultPath = os.Getenv("SystemRoot") + `\System32\drivers\etc\hosts`
