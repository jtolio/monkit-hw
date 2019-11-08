package hw

import (
	"bufio"
	"bytes"
	"flag"
	"os"

	"github.com/spacemonkeygo/monkit/v3"
)

var (
	oomLog = flag.String("monkit.hw.oomlog", "/var/log/kern.log",
		"path to log for oom notices")
)

func OOM() monkit.StatSource {
	// TODO: add oom scores
	return monkit.StatSourceFunc(
		func(cb func(key monkit.SeriesKey, field string, val float64)) {
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

			cb(monkit.NewSeriesKey("oom"), "count", float64(count))
		})
}

func init() { registrations = append(registrations, OOM()) }
