package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type PointMetrics struct {
	counter    prometheus.Counter
	errCounter prometheus.Counter
}

func NewPointMetrics() *PointMetrics {
	m := &PointMetrics{
		counter: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "points_created_total",
			Help: "Amount of created points",
		}),
		errCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "points_created_errors",
			Help: "Amount of errors wile creating points",
		}),
	}
	prometheus.MustRegister(m.counter)
	prometheus.MustRegister(m.errCounter)

	return m
}

func (m *PointMetrics) Inc() {
	m.counter.Inc()
}

func (m *PointMetrics) ErrInc() {
	m.errCounter.Inc()
}
