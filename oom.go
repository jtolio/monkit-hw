package hw

import (
	"bufio"
	"bytes"
	"flag"
	"os"

	"gopkg.in/spacemonkeygo/monkit.v2"
)

var (
	oomLog = flag.String("monkit.hw.oomlog", "/var/log/kern.log",
		"path to log for oom notices")
)

func OOM() monkit.StatSource {
	kills := monkit.Prefix("kills.", monkit.StatSourceFunc(
		func(cb func(name string, val float64)) {
			fh, err := os.Open(*oomLog)
			if err != nil {
				logger.Debuge(err)
				return
			}
			defer fh.Close()

			count := 0
			scanner := bufio.NewScanner(fh)
			for scanner.Scan() {
				if bytes.Contains(scanner.Bytes(), []byte("killed process")) {
					count++
				}
			}
			if err := scanner.Err(); err != nil {
				logger.Debuge(err)
				return
			}

			cb("total", float64(count))
		}))

	return monkit.StatSourceFunc(func(cb func(name string, val float64)) {
		kills.Stats(cb)
		// TODO: add oom scores
	})
}

func init() {
	registrations["oom"] = OOM()
}
