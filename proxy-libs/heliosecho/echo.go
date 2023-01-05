package heliosecho

import (
	origin_echo "github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

var InstrumentedSymbols = [...]string{"New"}

func addOpenTelemetryMiddleare(echo *origin_echo.Echo) {
	echo.Use(otelecho.Middleware("opentelemetry-middleware"))
}

func New() (e *origin_echo.Echo) {
	echo := origin_echo.New()
	addOpenTelemetryMiddleare(echo)
	return echo
}
