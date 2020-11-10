package env

import (
	"fmt"
	"runtime"
)

// Build time information.
type Build struct {
	App, Date, Commit string
}

func (b Build) String() string {
	return fmt.Sprintf("%s %s (%s revision %s) [%s-%s]", b.App, b.Date, b.Commit, runtime.GOARCH, runtime.GOOS, runtime.Version())
}

// Fields returns the build information as a map for structured logging.
func (b Build) Fields() map[string]interface{} {
	return map[string]interface{}{
		"app":        b.App,
		"build_date": b.Date,
		"commit":     b.Commit,
	}
}
