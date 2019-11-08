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
	key   monkit.SeriesKey
	field string
	val   float64
}

func IncludeDerivative(src monkit.StatSource) monkit.StatSource {
	var mtx sync.Mutex
	history := make([]captured, 0, derivWindows+1)
	return monkit.StatSourceFunc(func(cb func(key monkit.SeriesKey, field string, val float64)) {
		current := captured{seriesVals: map[string]seriesVal{}, ts: time.Now()}
		src.Stats(func(key monkit.SeriesKey, field string, val float64) {
			current.seriesVals[key.WithField(field)] = seriesVal{key, field, val}
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
				derivVal := (sVal.val - history[0].seriesVals[key].val) / timeDiff

				cb(sVal.key.WithTag("kind", "derivative"), sVal.field, derivVal)
				cb(sVal.key.WithTag("kind", "value"), sVal.field, sVal.val)
			}
		} else {
			for _, sVal := range current.seriesVals {
				cb(sVal.key.WithTag("kind", "value"), sVal.field, sVal.val)
			}
		}
	})
}
