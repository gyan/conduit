package dicontainer

import (
	"github.com/arpit32/conduit/api/controller"
	ca "github.com/yolobus/kuber/common/cadence"
	"github.com/yolobus/kuber/config"

	//ka "github.com/conduit/common/messaging"
	// db "github.com/conduit/common/mysql"
	"github.com/arpit32/conduit/api/service"
)

// ServiceContainer resolve all dependencies between controller, service, infrastructure except application level dependencies such us logging, config and etc ...
type ServiceContainer struct {
	config config.AppConfig

	//controllers
	PendulumController  *controller.PendulumController
	HTTPErrorController *controller.HTTPErrorController
}

// NewServiceContainer ...
func NewServiceContainer(config config.AppConfig) *ServiceContainer {
	return &ServiceContainer{
		config: config,
	}
}

//InitDependenciesInjection ...
func (container *ServiceContainer) InitDependenciesInjection() {
	//Initializing base controller
	baseController := controller.BaseController{Config: container.config}

	//Initializing Clients
	var cadenceClient ca.CadenceAdapter
	cadenceClient.Setup(&container.config.Cadence)

	// var mysqlClient db.DB
	// mysqlClient.ConnectSQL()

	//var kafkaClient ka.KafkaAdapter
	//kafkaClient.Setup(&container.config.Kafka)

	//Services
	pendulumService := &service.PendulumService{
		CadenceAdapter: cadenceClient,
		//KafkaAdapter: kafkaClient,
		// DB: mysqlClient,
		Logger: container.config.Logger,
	}

	//Initializing controllers
	container.HTTPErrorController = &controller.HTTPErrorController{BaseController: baseController}
	container.PendulumController = &controller.PendulumController{BaseController: baseController,
		PendulumService: pendulumService}

}
