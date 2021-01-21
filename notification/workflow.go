package notification

import (
	"context"

	"go.uber.org/cadence"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"

	"fmt"
	"time"
)

const notifyStaff = "startTripNotification"
const notifyCustomer = "notifyCustomerActivity"
const notifyCustomerAmenity = "notifyCustomerAmenityActivity"
const notifyCustomerWithTripStatus = "notifyCustomerWithTripStatusActivity"
const completeTrip = "completeTripActivity"

func init() {
	// workflow.RegisterWithOptions(
	// 	TripNotificationWorkflow,
	// 	workflow.RegisterOptions{Name: "tripNotificaiton"}
	// )

	activity.RegisterWithOptions(
		notifyStaffActivity,
		activity.RegisterOptions{Name: notifyStaff},
	)
	activity.RegisterWithOptions(
		notifyCustomerActivity,
		activity.RegisterOptions{Name: notifyCustomer},
	)
	activity.RegisterWithOptions(
		notifyCustomerAmenityActivity,
		activity.RegisterOptions{Name: notifyCustomerAmenity},
	)
	activity.RegisterWithOptions(
		notifyCustomerWithTripStatusActivity,
		activity.RegisterOptions{Name: notifyCustomerWithTripStatus},
	)
	activity.RegisterWithOptions(
		completeTripActivity,
		activity.RegisterOptions{Name: completeTrip},
	)
}

// TripNotificationWorkflow to notify passengers
func TripNotificationWorkflow(ctx workflow.Context, notificationsTasks []string, t Trip) error {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute * 3,
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

	for _, notificationTask := range notificationsTasks {
		switch task := notificationTask; task {
		case "notify_staff":
			err := workflow.ExecuteActivity(jobCtx, notifyStaff, t.Staff, t.Manager).Get(jobCtx, nil)
			if err != nil {
				logger.Error("Failed to execute notifyStaffActivity", zap.Error(err))
				return err
			}
		case "notify_customer":
			err := workflow.ExecuteActivity(jobCtx, notifyCustomer).Get(jobCtx, nil)
			if err != nil {
				logger.Error("Failed to execute notifyStaffActivity", zap.Error(err))
				return err
			}
		case "notify_customer_amenity":
			err := workflow.ExecuteActivity(jobCtx, notifyCustomerAmenity).Get(jobCtx, nil)
			if err != nil {
				logger.Error("Failed to execute notifyStaffActivity", zap.Error(err))
				return err
			}
		case "notify_customer_with_trip_status":
			err := workflow.ExecuteActivity(jobCtx, notifyCustomerWithTripStatus).Get(jobCtx, nil)
			if err != nil {
				logger.Error("Failed to execute notifyCustomerWithTripStatusActivity", zap.Error(err))
				return err
			}
		case "complete_trip":
			err := workflow.ExecuteActivity(jobCtx, completeTrip).Get(jobCtx, nil)
			if err != nil {
				logger.Error("Failed to execute notifyStaffActivity", zap.Error(err))
				return err
			}
		default:
			fmt.Println("Nothing to notify")
		}
	}
	// TODO : implement terminate instead of cancel
	_, cancel := workflow.WithCancel(ctx)
	cancel()
	return cadence.NewCanceledError("Hop Notified, Cancelling Child workflow ")
}

func notifyStaffActivity(ctx context.Context, staff, manager []string) error {
	// logger := activity.GetLogger(ctx)
	fmt.Println(time.Now(), "notifyStaffActivity  -> Start")
	fmt.Println(staff, manager)
	fmt.Println(time.Now(), "notifyStaffActivity -> Finished")

	return nil
}

func notifyCustomerActivity(ctx context.Context) error {
	// logger := activity.GetLogger(ctx)
	fmt.Println(time.Now(), "notifyCustomerActivity  -> Start")

	fmt.Println(time.Now(), "notifyCustomerActivity -> Finished")

	return nil
}

func notifyCustomerAmenityActivity(ctx context.Context) error {
	// logger := activity.GetLogger(ctx)
	fmt.Println(time.Now(), "notifyCustomerAmenityActivity  -> Start")

	fmt.Println(time.Now(), "notifyCustomerAmenityActivity -> Finished")

	return nil
}

func notifyCustomerWithTripStatusActivity(ctx context.Context) error {
	// logger := activity.GetLogger(ctx)
	fmt.Println(time.Now(), "notifyCustomerWithTripStatusActivity  -> Start")

	fmt.Println(time.Now(), "notifyCustomerWithTripStatusActivity -> Finished")

	return nil
}

func completeTripActivity(ctx context.Context) error {
	// logger := activity.GetLogger(ctx)
	fmt.Println(time.Now(), "completeTripActivity  -> Start")

	fmt.Println(time.Now(), "completeTripActivity -> Finished")

	return nil
}
