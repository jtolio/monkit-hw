// +build linux windows darwin,cgo

package hw

import (
	gosigar "github.com/cloudfoundry/gosigar"
	"github.com/spacemonkeygo/monkit/v3"
)

func Load() monkit.StatSource {
	return monkit.StatSourceFunc(func(cb func(key monkit.SeriesKey, field string, val float64)) {
		var load gosigar.LoadAverage
		err := load.Get()
		if err != nil {
			logger.Debuge(err)
			return
		}
		monkit.StatSourceFromStruct(monkit.NewSeriesKey("load"), &load).Stats(cb)
	})
}

func init() { registrations = append(registrations, Load()) }
