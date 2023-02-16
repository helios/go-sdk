package helioszerolog

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestZerolog(t *testing.T) {
	ctx := context.WithValue(context.Background(), "test", "test")
	newlog := log.With().Logger()
	newlog.WithContext(ctx)
	log.Ctx()

	newlog.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("goToHelios", "test2")
		})
	newlog.Info().Msg("hello world")
}
