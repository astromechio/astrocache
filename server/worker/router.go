package worker

import (
	"net/http"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/server/worker/handler"
	"github.com/gorilla/mux"
)

func router(app *config.App) *mux.Router {
	mux := mux.NewRouter()

	mux.Methods(http.MethodPost).Path("/v1/worker/block").HandlerFunc(handler.AddBlockHandler(app))

	mux.Methods(http.MethodPost).Path("/v1/value/{key}").HandlerFunc(handler.SetValueHandler(app))
	mux.Methods(http.MethodGet).Path("/v1/value/{key}").HandlerFunc(handler.GetValueHandler(app))

	return mux
}
