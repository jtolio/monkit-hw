// +build linux windows darwin,cgo

package hw

import (
	gosigar "github.com/cloudfoundry/gosigar"
	"github.com/spacemonkeygo/monkit/v3"
)

func CPU() monkit.StatSource {
	return IncludeDerivative(
		monkit.StatSourceFunc(func(cb func(key monkit.SeriesKey, field string, val float64)) {
			var cpu gosigar.Cpu
			err := cpu.Get()
			if err != nil {
				logger.Debuge(err)
				return
			}
			monkit.StatSourceFromStruct(monkit.NewSeriesKey("cpu"), &cpu).Stats(cb)
		}))
}

func init() { registrations = append(registrations, CPU()) }
