package pendulum

import (
	"fmt"
	"sort"
	"time"

	"go.uber.org/cadence"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"

	nt "github.com/arpit32/conduit/notification"
)

const (
	TaskList                          = "pendulum"
	SessionCreationErrorMsg           = "Session Creation Failed"
	startOscilliationActivityErrorMsg = "Failed to Start Oscillation"
	checkDriverStatusActivityErrorMsg = "Failed to check driver status"
	CallbackErrorEvent                = "error"
	Completed                         = "COMPLETED"
	ChildWorkflowExecErrMsg           = "Child Workflow execution failed"
)

func init() {
	workflow.RegisterWithOptions(
		Workflow,
		workflow.RegisterOptions{Name: TaskList},
	)

	workflow.RegisterWithOptions(
		nt.TripNotificationWorkflow,
		workflow.RegisterOptions{Name: "tripNotificaiton"},
	)
}

// Workflow to schedule child workflows
func Workflow(ctx workflow.Context, jobID string, q Query) error {

	logger := workflow.GetLogger(ctx)
	exec := workflow.GetInfo(ctx).WorkflowExecution

	runID := exec.RunID

	so := &workflow.SessionOptions{
		CreationTimeout:  time.Hour * 24,
		ExecutionTimeout: time.Minute * 10,
		HeartbeatTimeout: time.Minute * 3,
	}
	ctx, err := workflow.CreateSession(ctx, so)
	if err != nil {
		logger.Error(SessionCreationErrorMsg, zap.Error(err))
		return cadence.NewCustomError(err.Error(), SessionCreationErrorMsg)
	}
	defer workflow.CompleteSession(ctx)

	for q.CurrCityTask < len(q.Cities) {
		currCity := q.Cities[q.CurrCityTask]
		currCity.Status = "INIT"

		sort.Slice(currCity.Tasks, func(i, j int) bool {
			return currCity.Tasks[i].AlertMin < currCity.Tasks[j].AlertMin
		})

		cwo := workflow.ChildWorkflowOptions{
			WorkflowID:                   runID,
			TaskList:                     "pendulum",
			ExecutionStartToCloseTimeout: time.Hour * 24,
			TaskStartToCloseTimeout:      time.Minute * 24,
		}

		for _, task := range currCity.Tasks {
			cwo.CronSchedule = epochToCron(currCity.Etd + int64(task.AlertMin*60))
			childCtx := workflow.WithChildOptions(ctx, cwo)

			future := workflow.ExecuteChildWorkflow(childCtx, "tripNotificaiton", task.Name, q)
			currCity.Status = "SCHEDULED"
			err = future.Get(childCtx, nil)
			if err != nil {
				if cadence.IsCanceledError(err) {
					fmt.Println(task.Name, " ran for ", currCity.Name)
				} else {
					_, cancel := workflow.WithCancel(childCtx)
					cancel()
					return cadence.NewCanceledError(ChildWorkflowExecErrMsg, err.Error())
				}
			}
		}
		q.CurrCityTask++
	}

	return nil
}
