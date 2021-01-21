package api

import (
	"net/http"

	"github.com/arpit32/conduit/api/dicontainer"
	"github.com/arpit32/conduit/api/router"
	"github.com/yolobus/kuber/config"
)

//Application ...
type Application struct {
	config           config.AppConfig
	serviceContainer *dicontainer.ServiceContainer
	router           router.RoutingInterface
}

//New ...
func New(configPath string) *Application {
	var appConfig config.AppConfig
	appConfig.LoadConfig(configPath)

	return &Application{
		config: appConfig,
	}
}

// Init ...
func (app *Application) Init() {
	//start dependency injection
	app.serviceContainer = dicontainer.NewServiceContainer(app.config)

	app.serviceContainer.InitDependenciesInjection()

	//initialize new handlers
	app.router = router.NewRouter(app.config)
	app.router.Routes(app.serviceContainer)
}

// Start serve http server
func (app *Application) Start(serverPort string) error {
	return http.ListenAndServe(":"+serverPort, app.router.RouteMultiplexer())
}
