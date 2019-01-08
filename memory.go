// +build linux windows darwin,cgo

package hw

import (
	gosigar "github.com/cloudfoundry/gosigar"
	monkit "gopkg.in/spacemonkeygo/monkit.v2"
)

func Memory() monkit.StatSource {
	return monkit.StatSourceFunc(func(cb func(name string, val float64)) {
		var mem gosigar.Mem
		err := mem.Get()
		if err != nil {
			logger.Debuge(err)
			return
		}
		monkit.StatSourceFromStruct(&mem).Stats(cb)
		var swap gosigar.Swap
		err = swap.Get()
		if err != nil {
			logger.Debuge(err)
			return
		}
		monkit.Prefix("swap.", monkit.StatSourceFromStruct(&swap)).Stats(cb)
	})
}

func init() {
	registrations["memory"] = Memory()
}
