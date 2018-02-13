package common

import (
	"context"
	"net/http"

	"github.com/astromechio/astrocache/modes"
)

const (
	configKey = "astro.ctx.config"
)

// CtxWrap adds confg and keySet to a request context
func CtxWrap(config *modes.BaseConfig, inner http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxWithConfig(r, config)

		inner(w, r)
	}
}

func ctxWithConfig(r *http.Request, config *modes.BaseConfig) {
	r = r.WithContext(context.WithValue(r.Context(), configKey, config))
}

// ConfigFromReqCtx grabs a config from a request context
func ConfigFromReqCtx(r *http.Request) *modes.BaseConfig {
	config := r.Context().Value(configKey)

	return config.(*modes.BaseConfig)
}
