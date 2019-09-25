package hw

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/spacemonkeygo/monkit/v3"
)

func statSourceFromDir(dir string) monkit.StatSource {
	return monkit.StatSourceFunc(
		func(cb func(series monkit.Series, val float64)) {
			entries, err := ioutil.ReadDir(dir)
			if err != nil {
				logger.Debuge(err)
				return
			}
			for _, entry := range entries {
				data, err := ioutil.ReadFile(filepath.Join(dir, entry.Name()))
				if err != nil {
					logger.Debuge(err)
					continue
				}
				val, err := strconv.ParseFloat(string(bytes.TrimSpace(data)), 64)
				if err != nil {
					logger.Debuge(err)
					continue
				}
				cb(monkit.NewSeries("hardware", entry.Name()), val)
			}
		})
}
