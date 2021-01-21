package router

import (
	"github.com/arpit32/conduit/api/dicontainer"
	"github.com/go-chi/cors"
	"github.com/yolobus/kuber/config"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

//RoutingInterface ...
type RoutingInterface interface {
	Routes(serviceContainer *dicontainer.ServiceContainer)
	RouteMultiplexer() *chi.Mux
}

type router struct {
	config config.AppConfig
	mux    *chi.Mux
}

//NewRouter ...
func NewRouter(generalConfig config.AppConfig) RoutingInterface {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.RealIP)
	// mux.Use(SetJSON)
	// mux.Use(logger.NewStructuredLogger())
	mux.Use(cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}).Handler)
	return &router{
		mux:    mux,
		config: generalConfig,
		// logger: logger,
	}
}

func (h *router) RouteMultiplexer() *chi.Mux {
	return h.mux
}

func (h *router) Routes(container *dicontainer.ServiceContainer) {
	h.mux.Group(func(r chi.Router) {

		r.Post("/v1/pendulum/start", container.PendulumController.CreateJob)
	})

	h.mux.NotFound(container.HTTPErrorController.ResourceNotFound)
}
