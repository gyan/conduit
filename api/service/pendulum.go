package service

import (
	"context"
	"strings"
	"time"

	m "github.com/arpit32/conduit/api/model"
	pm "github.com/arpit32/conduit/pendulum"
	"github.com/google/uuid"

	// "github.com/google/uuid"
	ca "github.com/yolobus/kuber/common/cadence"

	// ka "github.com/conduit/common/messaging"

	// "go.uber.org/cadence/client"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

//PendulumService ...
type PendulumService struct {
	CadenceAdapter ca.CadenceAdapter
	// KafkaAdapter ka.KafkaAdapter
	// DB 	mysql.DB
	Logger *zap.Logger
}

// CreateJob ...
func (b *PendulumService) CreateJob(ctx context.Context, trip m.Trip) (*workflow.Execution, error) {

	if strings.ToLower(trip.IsDryRun) == "true" {
		for _, c := range trip.Cities {
			c.Etd = time.Now().Unix()
		}
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:                              "pendulum_" + uuid.New().String(),
		TaskList:                        pm.TaskList,
		ExecutionStartToCloseTimeout:    time.Hour * 24,
		DecisionTaskStartToCloseTimeout: time.Minute * 24,
	}

	execution, err := b.CadenceAdapter.CadenceClient.StartWorkflow(
		context.Background(),
		workflowOptions,
		pm.TaskList,
		uuid.New().String(),
		trip,
	)

	return execution, err
}
