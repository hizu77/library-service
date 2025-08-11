package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const metricsPath = "/metrics"

func InitMetrics(app fiber.Router) {
	app.Get(
		metricsPath,
		adaptor.HTTPHandler(promhttp.Handler()),
	)
}
