// +build !freebsd !darwin darwin,cgo

package hw

import (
	gosigar "github.com/cloudfoundry/gosigar"
	monkit "gopkg.in/spacemonkeygo/monkit.v2"
)

func Load() monkit.StatSource {
	return monkit.StatSourceFunc(func(cb func(name string, val float64)) {
		var load gosigar.LoadAverage
		err := load.Get()
		if err != nil {
			logger.Debuge(err)
			return
		}
		monkit.StatSourceFromStruct(&load).Stats(cb)
	})
}

func init() {
	registrations["load"] = Load()
}
