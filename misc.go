// +build linux windows darwin,cgo

package hw

import (
	gosigar "github.com/cloudfoundry/gosigar"
	"github.com/spacemonkeygo/monkit/v3"
)

// uptime, control
func Misc() monkit.StatSource {
	return monkit.StatSourceFunc(func(cb func(series monkit.Series, val float64)) {
		cb(monkit.NewSeries("control", "value"), 1)
		var u gosigar.Uptime
		err := u.Get()
		if err != nil {
			logger.Debuge(err)
			return
		}
		cb(monkit.NewSeries("uptime", "value"), u.Length)
	})
}

func init() {
	registrations["misc"] = Misc()
}
