// +build linux

package hw

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/spacemonkeygo/monkit.v2"
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
	return monkit.StatSourceFunc(func(cb func(name string, val float64)) {
		v4conns, err := countConns("/proc/net/tcp")
		if err != nil {
			logger.Debuge(err)
			return
		}
		cb("v4", float64(v4conns))

		v6conns, err := countConns("/proc/net/tcp6")
		if err != nil {
			logger.Debuge(err)
			return
		}
		cb("v6", float64(v6conns))
		cb("total", float64(v4conns+v6conns))
	})
}

func NetStats() monkit.StatSource {
	return IncludeDerivative(
		monkit.StatSourceFunc(func(cb func(name string, val float64)) {
			interfaces, err := ioutil.ReadDir("/sys/class/net")
			if err != nil {
				logger.Debuge(err)
				return
			}
			for _, iface := range interfaces {
				statsdir := filepath.Join("/sys/class/net", iface.Name(), "statistics")
				monkit.Prefix(iface.Name()+".", statSourceFromDir(statsdir)).Stats(cb)
			}
		}))
}

func Network() monkit.StatSource {
	conns := monkit.Prefix("conns.", Conns())
	stats := monkit.Prefix("stats.", NetStats())
	return monkit.StatSourceFunc(func(cb func(name string, val float64)) {
		conns.Stats(cb)
		stats.Stats(cb)
	})
}

func init() {
	registrations["network"] = Network()
}
