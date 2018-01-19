package hw

import (
	"github.com/spacemonkeygo/spacelog"
	"gopkg.in/spacemonkeygo/monkit.v2"
)

var (
	registrations = map[string]monkit.StatSource{}
	logger        = spacelog.GetLogger()
)

func Register(registry *monkit.Registry) {
	if registry == nil {
		registry = monkit.Default
	}
	pkg := registry.ScopeNamed("hw")
	for name, source := range registrations {
		pkg.Chain(name, source)
	}
}
