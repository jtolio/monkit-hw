package hw

import (
	"sync"
	"time"

	"github.com/spacemonkeygo/monkit/v3"
)

const (
	derivWindows = 4
)

type captured struct {
	seriesVals map[string]seriesVal
	ts         time.Time
}

type seriesVal struct {
	series monkit.Series
	val    float64
}

func IncludeDerivative(src monkit.StatSource) monkit.StatSource {
	var mtx sync.Mutex
	history := make([]captured, 0, derivWindows+1)
	return monkit.StatSourceFunc(func(cb func(series monkit.Series, val float64)) {
		current := captured{seriesVals: map[string]seriesVal{}, ts: time.Now()}
		src.Stats(func(series monkit.Series, val float64) {
			current.seriesVals[series.String()] = seriesVal{series, val}
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
			for key, sVal := range current.seriesVals {
				derivVal := (sVal.val-history[0].seriesVals[key].val)/timeDiff
				cb(monkit.NewSeries("hardware", key+".deriv"), derivVal)
				cb(monkit.NewSeries("hardware", key+".val"), sVal.val)
			}
		} else {
			for key, sVal := range current.seriesVals {
				cb(monkit.NewSeries("hardware", key+".seriesVal"), sVal.val)
			}
		}
	})
}
