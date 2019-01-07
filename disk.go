// +build !freebsd !darwin darwin,cgo

package hw

import (
	"strings"

	gosigar "github.com/cloudfoundry/gosigar"
	monkit "gopkg.in/spacemonkeygo/monkit.v2"
)

func Disk() monkit.StatSource {
	return monkit.StatSourceFunc(func(cb func(name string, val float64)) {
		var fslist gosigar.FileSystemList
		err := fslist.Get()
		if err != nil {
			logger.Debuge(err)
			return
		}
		for _, fs := range fslist.List {
			if !strings.HasPrefix(fs.DevName, "/") {
				continue
			}
			var fsu gosigar.FileSystemUsage
			err = fsu.Get(fs.DirName)
			if err != nil {
				logger.Debuge(err)
				continue
			}
			monkit.Prefix(fs.DevName+".", monkit.StatSourceFromStruct(&fsu)).Stats(cb)
		}
	})
}

func init() {
	registrations["disk"] = Disk()
}
