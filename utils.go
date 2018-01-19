package hw

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"gopkg.in/spacemonkeygo/monkit.v2"
)

func statSourceFromDir(dir string) monkit.StatSource {
	return monkit.StatSourceFunc(
		func(cb func(name string, val float64)) {
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
				cb(entry.Name(), val)
			}
		})
}
