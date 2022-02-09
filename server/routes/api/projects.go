package api

import "github.com/go-chi/chi/v5"

func newProjectsApiRouter() chi.Router {
	r := chi.NewRouter()

	return r
}

func newSubprojectsApiRouter() chi.Router {
	r := chi.NewRouter()

	return r
}

func newProjectAclRouter() chi.Router {
	r := chi.NewRouter()

	return r
}
