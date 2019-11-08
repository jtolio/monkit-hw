// +build linux windows darwin,cgo

package hw

import (
	"strings"

	gosigar "github.com/cloudfoundry/gosigar"
	"github.com/spacemonkeygo/monkit/v3"
)

func Disk() monkit.StatSource {
	return monkit.StatSourceFunc(func(cb func(key monkit.SeriesKey, field string, val float64)) {
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
			monkit.StatSourceFromStruct(monkit.NewSeriesKey("disk"), &fsu).Stats(
				func(key monkit.SeriesKey, field string, val float64) {
					cb(key.WithTag("device", fs.DevName), field, val)
				})
		}
	})
}

func init() { registrations = append(registrations, Disk()) }
