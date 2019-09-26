// +build linux windows darwin,cgo

package hw

import (
	gosigar "github.com/cloudfoundry/gosigar"
	monkit "github.com/spacemonkeygo/monkit/v3"
)

func Load() monkit.StatSource {
	return monkit.StatSourceFunc(func(cb func(series monkit.Series, val float64)) {
		var load gosigar.LoadAverage
		err := load.Get()
		if err != nil {
			logger.Debuge(err)
			return
		}
		monkit.StatSourceFromStruct(&load).Stats(func(series monkit.Series, val float64) {
			series.Measurement = "hardware"
			series.Tags = series.Tags.Set("kind", "load")
			cb(series, val)
		})
	})
}

func init() {
	registrations["load"] = Load()
}
