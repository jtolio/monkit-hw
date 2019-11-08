package hw

import (
	"github.com/spacemonkeygo/monkit/v3"
	"github.com/spacemonkeygo/spacelog"
)

var (
	registrations = []monkit.StatSource{}
	logger        = spacelog.GetLogger()
)

func Register(registry *monkit.Registry) {
	if registry == nil {
		registry = monkit.Default
	}
	pkg := registry.Package()
	for _, source := range registrations {
		pkg.Chain(source)
	}
}
