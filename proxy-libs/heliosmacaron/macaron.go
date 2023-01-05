package heliosmacaron

import (
	"io"

	"go.opentelemetry.io/contrib/instrumentation/gopkg.in/macaron.v1/otelmacaron"
	origin_macaron "gopkg.in/macaron.v1"
)

var InstrumentedSymbols = [...]string{"NewWithLogger", "Classic", "New"}

func addOtelMiddleware(mac *origin_macaron.Macaron) {
	mac.Use(otelmacaron.Middleware("opentelemetry-middleware"))
}

func NewWithLogger(out io.Writer) *origin_macaron.Macaron {
	mac := origin_macaron.NewWithLogger(out)
	addOtelMiddleware(mac)
	return mac
}

func New() *origin_macaron.Macaron {
	mac := origin_macaron.New()
	addOtelMiddleware(mac)
	return mac
}

func Classic() *origin_macaron.Macaron {
	mac := origin_macaron.Classic()
	addOtelMiddleware(mac)
	return mac
}
