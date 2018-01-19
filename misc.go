package hw

import (
	gosigar "github.com/cloudfoundry/gosigar"
	"gopkg.in/spacemonkeygo/monkit.v2"
)

// uptime, control
func Misc() monkit.StatSource {
	return monkit.StatSourceFunc(func(cb func(name string, val float64)) {
		cb("control", 1)
		var u gosigar.Uptime
		err := u.Get()
		if err != nil {
			logger.Debuge(err)
			return
		}
		cb("uptime", u.Length)
	})
}

func init() {
	registrations["misc"] = Misc()
}
