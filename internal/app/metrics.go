package app

import (
	"github.com/4udiwe/avito-pvz/internal/metrics"
)

func (app *App) PointMetrics() *metrics.PointMetrics {
	if app.pointMetrics != nil {
		return app.pointMetrics
	}
	app.pointMetrics = metrics.NewPointMetrics()
	return app.pointMetrics
}

func (app *App) ProductMetrics() *metrics.ProductMetrics {
	if app.productMetrics != nil {
		return app.productMetrics
	}
	app.productMetrics = metrics.NewProductMetrics()
	return app.productMetrics
}

func (app *App) ReceptionMetrics() *metrics.ReceptionMetrics {
	if app.receptionMetrics != nil {
		return app.receptionMetrics
	}
	app.receptionMetrics = metrics.NewReceptionMetrics()
	return app.receptionMetrics
}
