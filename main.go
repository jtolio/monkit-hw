package hw

import (
	"github.com/spacemonkeygo/monkit/v3"
	"github.com/spacemonkeygo/spacelog"
)

var (
	registrations = map[string]monkit.StatSource{}
	logger        = spacelog.GetLogger()
)

func Register(registry *monkit.Registry) {
	if registry == nil {
		registry = monkit.Default
	}
	pkg := registry.Package()
	for name, source := range registrations {
		pkg.Chain(name, source)
	}
}
