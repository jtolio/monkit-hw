// +build linux

package hw

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spacemonkeygo/monkit/v3"
)

func countConns(path string) (amount int, err error) {
	fh, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer fh.Close()

	scanner := bufio.NewScanner(fh)
	scanner.Scan() // skip the header
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) < 4 {
			continue
		}
		state := line[3]
		if state != "0A" {
			amount += 1
		}
	}
	err = scanner.Err()
	if err != nil {
		return 0, err
	}
	return amount, nil
}

func Conns() monkit.StatSource {
	return monkit.StatSourceFunc(func(cb func(series monkit.Series, val float64)) {
		series := monkit.NewSeries("conns", "connections")

		v4conns, err := countConns("/proc/net/tcp")
		if err != nil {
			logger.Debuge(err)
			return
		}
		series.Tags = series.Tags.Set("protocol", "ipv4")
		cb(series, float64(v4conns))

		v6conns, err := countConns("/proc/net/tcp6")
		if err != nil {
			logger.Debuge(err)
			return
		}

		series.Tags = series.Tags.Set("protocol", "ipv6")
		cb(series, float64(v6conns))

		series.Tags = series.Tags.Set("protocol", "all")
		cb(series, float64(v4conns+v6conns))
	})
}

func NetStats() monkit.StatSource {
	return IncludeDerivative(
		monkit.StatSourceFunc(func(cb func(series monkit.Series, val float64)) {
			interfaces, err := ioutil.ReadDir("/sys/class/net")
			if err != nil {
				logger.Debuge(err)
				return
			}
			for _, iface := range interfaces {
				statsdir := filepath.Join("/sys/class/net", iface.Name(), "statistics")
				statSourceFromDir(statsdir).Stats(func(series monkit.Series, val float64) {
					series.Measurement = "netstats"
					series.Tags = series.Tags.Set("interface", iface.Name())
					cb(series, val)
				})
			}
		}))
}

func Network() monkit.StatSource {
	return monkit.StatSourceFunc(func(cb func(series monkit.Series, val float64)) {
		Conns().Stats(cb)
		NetStats().Stats(cb)
	})
}

func init() {
	registrations["network"] = Network()
}
