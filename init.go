package sm

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	db      Store
	metrics *Metrics
)

func InitMetrics(name string) *Metrics {
	metrics = NewMetrics(name)

	ctime := time.Now()

	// server uptime counter
	metrics.NewCounterFunc(
		"server", "uptime",
		"Number of nanoseconds the server has been running",
		func() float64 {
			return float64(time.Since(ctime).Nanoseconds())
		},
	)

	// server requests counter
	metrics.NewCounterVec(
		"server", "requests",
		"Number of requests made to the server",
		[]string{"method", "path"},
	)

	// event count counter
	metrics.NewCounter(
		"event", "count",
		"Job count",
	)

	// event index summary
	metrics.NewSummary(
		"event", "index",
		"Index duration in seconds",
	)

	return metrics
}

func InitDB(uri string) (Store, error) {
	u, err := ParseURI(uri)
	if err != nil {
		log.Errorf("error parsing db uri %s: %s", uri, err)
		return nil, err
	}

	switch u.Type {
	case "memory":
		db, err = NewMemoryStore()
		if err != nil {
			log.Errorf("error creating store %s: %s", uri, err)
			return nil, err
		}
		log.Infof("Using MemoryStore %s", uri)
		return db, nil
	case "bolt":
		db, err = NewBoltStore(u.Path)
		if err != nil {
			log.Errorf("error creating store %s: %s", uri, err)
			return nil, err
		}
		log.Infof("Using BoltStore %s", uri)
		return db, nil
	default:
		err := fmt.Errorf("unsupported db uri: %s", uri)
		log.Error(err)
		return nil, err
	}
}
