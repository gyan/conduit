package child

import (
	"context"

	"go.uber.org/cadence"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"

	"fmt"
	"time"
)

const startTripNotification = "startTripNotification"

func init() {
	// workflow.RegisterWithOptions(TripNotificationWorkflow, workflow.RegisterOptions{Name: "tripNotificaiton"})

	activity.RegisterWithOptions(
		startTripNotificationActivity,
		activity.RegisterOptions{Name: startTripNotification},
	)
}

// TripNotificationWorkflow to notify passengers
func TripNotificationWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Second * 20,
		HeartbeatTimeout:       time.Second * 10,
		RetryPolicy: &cadence.RetryPolicy{
			InitialInterval:          time.Second,
			BackoffCoefficient:       2.0,
			MaximumInterval:          time.Minute * 5,
			ExpirationInterval:       time.Hour * 10,
			MaximumAttempts:          2,
			NonRetriableErrorReasons: []string{"bad-error"},
		},
	}
	jobCtx := workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(jobCtx)

	var notifyHop = false
	err := workflow.ExecuteActivity(jobCtx, startTripNotificationActivity).Get(jobCtx, nil)
	if err != nil {
		logger.Error("Failed to execute tripNotificaitonActivity function", zap.Error(err))
		return err
	}
	notifyHop = true
	fmt.Println("Hop notified: ", notifyHop)

	_, cancel := workflow.WithCancel(ctx)
	cancel()
	return cadence.NewCanceledError("Hop Notified, Cancelling Child workflow ")
}

func startTripNotificationActivity(ctx context.Context) error {
	// logger := activity.GetLogger(ctx)
	fmt.Println(time.Now(), "tripNotificaitonActivity  -> Start")

	fmt.Println(time.Now(), "tripNotificaitonActivity -> Finished")

	return nil
}
