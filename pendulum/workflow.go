package pendulum

import (
	"fmt"
	"time"

	"go.uber.org/cadence"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"

	ch "github.com/arpit32/conduit/pendulum/child"
)

const (
	TaskList                          = "pendulum"
	SessionCreationErrorMsg           = "Session Creation Failed"
	startOscilliationActivityErrorMsg = "Failed to Start Oscillation"
	checkDriverStatusActivityErrorMsg = "Failed to check driver status"
	Task                              = "task"
	CallbackErrorEvent                = "error"
	Completed                         = "COMPLETED"
	ChildWorkflowExecErrMsg           = "Child Workflow execution failed"
)

func init() {
	workflow.RegisterWithOptions(Workflow, workflow.RegisterOptions{Name: TaskList})
	workflow.RegisterWithOptions(ch.TripNotificationWorkflow, workflow.RegisterOptions{Name: "tripNotificaiton"})
}

// [10, 15, 24]

// Workflow to schedule child workflows
func Workflow(ctx workflow.Context) error {

	logger := workflow.GetLogger(ctx)
	exec := workflow.GetInfo(ctx).WorkflowExecution

	runID := exec.RunID

	so := &workflow.SessionOptions{
		CreationTimeout:  time.Hour * 24,
		ExecutionTimeout: time.Minute * 5,
		HeartbeatTimeout: time.Minute * 3,
	}
	ctx, err := workflow.CreateSession(ctx, so)
	if err != nil {
		logger.Error(SessionCreationErrorMsg, zap.Error(err))
		return cadence.NewCustomError(err.Error(), SessionCreationErrorMsg)
	}
	defer workflow.CompleteSession(ctx)
	
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute * 5,
		ScheduleToCloseTimeout: time.Minute * 5,
		HeartbeatTimeout:       time.Minute * 3,
		RetryPolicy: &cadence.RetryPolicy{
			InitialInterval:          time.Second,
			BackoffCoefficient:       2.0,
			MaximumInterval:          time.Minute * 5,
			ExpirationInterval:       time.Hour * 10,
			MaximumAttempts:          2,
			NonRetriableErrorReasons: []string{"bad-error"},
		},
	}
	activityCtx := workflow.WithActivityOptions(ctx, ao)

	// trip starts at : 11:37
	// notify first hop at 11:35
	// check driver status
	// Update trip status at 11:38

	// ETA to next hop : 11:39
	// Notify deboarded passengers

	hops := [3]int{1, 1, 1}
	var driverStatus string
	var tripStarted bool

	for ix, hop := range hops {
		cwo := workflow.ChildWorkflowOptions{
			WorkflowID:                   runID,
			TaskList:                     "pendulum",
			ExecutionStartToCloseTimeout: time.Hour * 24,
			TaskStartToCloseTimeout:      time.Minute * 24,
			CronSchedule:                 fmt.Sprintf("*/%v * * * *", hop),
		}

		childCtx := workflow.WithChildOptions(ctx, cwo)

		// Child workflow sch3duled to notify first hop 10 seconds before the trip starts
		err = workflow.ExecuteChildWorkflow(childCtx, "tripNotificaiton", runID).Get(childCtx, nil)
		if err != nil {
			if cadence.IsCanceledError(err) {
				fmt.Println("Child Workflow ran for hop", ix+1, "at", hop, " second")
			} else {
				_, cancel := workflow.WithCancel(childCtx)
				cancel()
				return cadence.NewCanceledError(ChildWorkflowExecErrMsg, err.Error())
			}
		}

		// Tasks to do before the trip starts
		if ix == 0 {
			err = workflow.ExecuteActivity(activityCtx, checkDriverStatus).Get(ctx, &driverStatus)
			if err != nil {
				logger.Error(checkDriverStatusActivityErrorMsg, zap.Error(err))
				return cadence.NewCustomError(err.Error(), checkDriverStatusActivityErrorMsg)
			}
			if driverStatus == "ready" {
				fmt.Println("Driver Status:", driverStatus)
				workflow.Sleep(ctx, time.Second*3)
				tripStarted = true
				fmt.Println("Trip Started: ", tripStarted)
			}
		}
	}

	return nil
}
