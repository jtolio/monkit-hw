package hw

import (
	gosigar "github.com/cloudfoundry/gosigar"
	"gopkg.in/spacemonkeygo/monkit.v2"
)

func CPU() monkit.StatSource {
	return IncludeDerivative(
		monkit.StatSourceFunc(func(cb func(name string, val float64)) {
			var cpu gosigar.Cpu
			err := cpu.Get()
			if err != nil {
				logger.Debuge(err)
				return
			}
			monkit.StatSourceFromStruct(&cpu).Stats(cb)
		}))
}

func init() {
	registrations["cpu"] = CPU()
}
