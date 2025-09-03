package metrics

import "github.com/prometheus/client_golang/prometheus"

type ReceptionMetrics struct {
	counter    prometheus.Counter
	errCounter prometheus.Counter
}

func NewReceptionMetrics() *ReceptionMetrics {
	m := &ReceptionMetrics{
		counter: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "reception_created_total",
			Help: "Amount of openned receptions",
		}),
		errCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "reception_created_errors",
			Help: "Amount of errors wile openning reception",
		}),
	}

	prometheus.MustRegister(m.counter)
	prometheus.MustRegister(m.errCounter)

	return m
}

func (m *ReceptionMetrics) Inc() {
	m.counter.Inc()
}

func (m *ReceptionMetrics) ErrInc() {
	m.errCounter.Inc()
}
