// +build linux windows darwin,cgo

package hw

import (
	gosigar "github.com/cloudfoundry/gosigar"
	"github.com/spacemonkeygo/monkit/v3"
)

func Memory() monkit.StatSource {
	return monkit.StatSourceFunc(func(cb func(key monkit.SeriesKey, field string, val float64)) {
		var mem gosigar.Mem
		err := mem.Get()
		if err != nil {
			logger.Debuge(err)
			return
		}
		monkit.StatSourceFromStruct(monkit.NewSeriesKey("memory"), &mem).Stats(cb)

		var swap gosigar.Swap
		err = swap.Get()
		if err != nil {
			logger.Debuge(err)
			return
		}
		monkit.StatSourceFromStruct(monkit.NewSeriesKey("swap"), &swap).Stats(cb)
	})
}

func init() { registrations = append(registrations, Memory()) }
