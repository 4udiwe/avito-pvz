package metrics

import "github.com/prometheus/client_golang/prometheus"

type ProductMetrics struct {
	counter    prometheus.Counter
	errCounter prometheus.Counter
}

func NewProductMetrics() *ProductMetrics {
	m := &ProductMetrics{
		counter: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "products_created_total",
			Help: "Amount of created products",
		}),
		errCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "products_created_errors",
			Help: "Amount of errors wile creating products",
		}),
	}

	prometheus.MustRegister(m.counter)
	prometheus.MustRegister(m.errCounter)

	return m
}

func (m *ProductMetrics) Inc() {
	m.counter.Inc()
}

func (m *ProductMetrics) ErrInc() {
	m.errCounter.Inc()
}
