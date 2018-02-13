package master

import "github.com/gorilla/mux"

func router(config *Config) *mux.Router {
	mux := mux.NewRouter()

	return mux
}
