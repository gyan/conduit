package pendulum // called the package pendulum cause "intervals"

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/yolobus/kuber/config"
	"go.uber.org/cadence/worker"
	"go.uber.org/zap"

	_ "github.com/joho/godotenv/autoload"

	ca "github.com/yolobus/kuber/common/cadence"
)

//Worker ...
type Worker struct {
	config         config.AppConfig
	taskList       string
	cadenceAdapter ca.CadenceAdapter
	// kafkaAdapter   ka.KafkaAdapter
	options worker.Options
	logger  log.Logger
}

//New ...
func New(configPath string) *Worker {
	var appConfig config.AppConfig
	appConfig.LoadConfig(configPath)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "upload",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	return &Worker{
		config: appConfig,
		logger: logger,
	}
}

// Init ...
func (w *Worker) Init(tasklist, verbose, workerType string) {
	//start dependency injection
	w.cadenceAdapter.Setup(&w.config.Cadence)
	workerOptions := worker.Options{
		MetricsScope:          w.cadenceAdapter.Scope,
		EnableLoggingInReplay: true,
	}

	if workerType == "activity" {
		fmt.Println("activity")

		ctx := context.Background()
		ctx = context.WithValue(ctx, "cadenceClient", w.cadenceAdapter)

		workerOptions.BackgroundActivityContext = ctx
		workerOptions.EnableSessionWorker = true
		workerOptions.DisableWorkflowWorker = true
		workerOptions.DisableActivityWorker = false
		workerOptions.MaxConcurrentSessionExecutionSize = 10
		workerOptions.WorkerStopTimeout = time.Second * 10

	} else if workerType == "workflow" {
		workerOptions.EnableSessionWorker = false
		workerOptions.DisableWorkflowWorker = false
		workerOptions.DisableActivityWorker = true
		workerOptions.WorkerStopTimeout = time.Second * 10
	}
	if verbose == "0" {
		workerOptions.Logger = zap.NewNop()
	} else {
		workerOptions.Logger = w.cadenceAdapter.Logger
	}
	w.options = workerOptions
	w.taskList = tasklist
}

//Start ...
func (w *Worker) Start() {
	// Configure worker options.

	cadenceWorker := worker.New(w.cadenceAdapter.ServiceClient, w.config.Cadence.Domain, w.taskList, w.options)
	err := cadenceWorker.Start()
	if err != nil {
		w.cadenceAdapter.Logger.Error("Failed to start workers.", zap.Error(err))
		panic("Failed to start workers")
	}

	done := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		killSignal := <-sigint
		switch killSignal {
		case os.Interrupt:
			level.Info(w.logger).Log("Got SIGINT...")
		case syscall.SIGTERM:
			level.Info(w.logger).Log("Got SIGTERM...")
		}
		time.Sleep(time.Second * 2)
		cadenceWorker.Stop()
		close(done)
	}()
	<-done
}
