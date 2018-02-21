package hw

import (
	"sync"
	"time"

	"gopkg.in/spacemonkeygo/monkit.v2"
)

const (
	derivWindows = 4
)

type captured struct {
	vals map[string]float64
	ts   time.Time
}

func IncludeDerivative(src monkit.StatSource) monkit.StatSource {
	var mtx sync.Mutex
	history := make([]captured, 0, derivWindows+1)
	return monkit.StatSourceFunc(func(cb func(name string, val float64)) {
		current := captured{vals: map[string]float64{}, ts: time.Now()}
		src.Stats(func(name string, val float64) {
			current.vals[name] = val
		})
		mtx.Lock()
		defer mtx.Unlock()
		history = append(history, current)
		if len(history) > derivWindows {
			copy(history, history[1:])
			history = history[:derivWindows]
		}
		timeDiff := current.ts.Sub(history[0].ts).Seconds()
		if timeDiff > 0 {
			for key, val := range current.vals {
				cb(key+".deriv", (val-history[0].vals[key])/timeDiff)
				cb(key+".val", val)
			}
		} else {
			for key, val := range current.vals {
				cb(key+".val", val)
			}
		}
	})
}
