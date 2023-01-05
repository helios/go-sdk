package heliosgin

import (
	"net/http"

	origin_gin "github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var InstrumentedSymbols = [...]string{"New", "Default", "CreateTestContext"}

func addOpenTelemetryMiddleware(engine *origin_gin.Engine) {
	engine.Use(otelgin.Middleware("opentelemetry-middleware"))
}

func CreateTestContext(w http.ResponseWriter) (c *origin_gin.Context, r *origin_gin.Engine) {
	ctx, engine := origin_gin.CreateTestContext(w)
	addOpenTelemetryMiddleware(engine)
	return ctx, engine
}

func New() *origin_gin.Engine {
	engine := origin_gin.New()
	addOpenTelemetryMiddleware(engine)
	return engine
}

func Default() *origin_gin.Engine {
	engine := origin_gin.Default()
	addOpenTelemetryMiddleware(engine)
	return engine
}
