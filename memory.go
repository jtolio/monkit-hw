// +build linux windows darwin,cgo

package hw

import (
	gosigar "github.com/cloudfoundry/gosigar"
	monkit "github.com/spacemonkeygo/monkit/v3"
)

func Memory() monkit.StatSource {
	return monkit.StatSourceFunc(func(cb func(series monkit.Series, val float64)) {
		var mem gosigar.Mem
		err := mem.Get()
		if err != nil {
			logger.Debuge(err)
			return
		}
		monkit.StatSourceFromStruct(&mem).Stats(func(series monkit.Series, val float64) {
			series.Measurement = "hardware"
			series.Tags = series.Tags.Set("kind", "memory")
			cb(series, val)
		})
		var swap gosigar.Swap
		err = swap.Get()
		if err != nil {
			logger.Debuge(err)
			return
		}
		monkit.StatSourceFromStruct(&swap).Stats(func(series monkit.Series, val float64) {
			series.Measurement = "hardware"
			series.Tags = series.Tags.Set("kind", "swap")
			cb(series, val)
		})
	})
}

func init() {
	registrations["memory"] = Memory()
}
