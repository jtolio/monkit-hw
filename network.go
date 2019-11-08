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
	return monkit.StatSourceFunc(func(cb func(key monkit.SeriesKey, field string, val float64)) {
		v4conns, err := countConns("/proc/net/tcp")
		if err != nil {
			logger.Debuge(err)
			return
		}
		cb(monkit.NewSeriesKey("conns").WithTag("protocol", "ipv4"), "count", float64(v4conns))

		v6conns, err := countConns("/proc/net/tcp6")
		if err != nil {
			logger.Debuge(err)
			return
		}

		cb(monkit.NewSeriesKey("conns").WithTag("protocol", "ipv6"), "count", float64(v6conns))
		cb(monkit.NewSeriesKey("conns").WithTag("protocol", "all"), "count", float64(v4conns+v6conns))
	})
}

func NetStats() monkit.StatSource {
	return IncludeDerivative(
		monkit.StatSourceFunc(func(cb func(key monkit.SeriesKey, field string, val float64)) {
			interfaces, err := ioutil.ReadDir("/sys/class/net")
			if err != nil {
				logger.Debuge(err)
				return
			}
			for _, iface := range interfaces {
				statsdir := filepath.Join("/sys/class/net", iface.Name(), "statistics")
				statSourceFromDir("netstats", statsdir).Stats(func(key monkit.SeriesKey, field string, val float64) {
					cb(key.WithTag("interface", iface.Name()), field, val)
				})
			}
		}))
}

func Network() monkit.StatSource {
	return monkit.StatSourceFunc(func(cb func(key monkit.SeriesKey, field string, val float64)) {
		Conns().Stats(cb)
		NetStats().Stats(cb)
	})
}

func init() { registrations = append(registrations, Network()) }
