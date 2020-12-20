package pendulum

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/cadence/activity"
)

/**
 * Sample activities used by Pendulum workflow.
 */
const (
	startOscillation  = "startOscillation"
	checkDriverStatus = "checkDriverStatus"
)

// This is registration process where you register all your activity handlers.
func init() {
	activity.RegisterWithOptions(
		startOscillationActivity,
		activity.RegisterOptions{Name: startOscillation},
	)
	activity.RegisterWithOptions(
		checkDriverStatusActivity,
		activity.RegisterOptions{Name: checkDriverStatus},
	)
}

func startOscillationActivity(ctx context.Context) error {
	// logger := activity.GetLogger(ctx)

	fmt.Println(time.Now(), "oscillation activity -> Start")

	fmt.Println(time.Now(), "oscillation activity -> Finished")

	return nil
}

func checkDriverStatusActivity(ctx context.Context) (string, error) {
	// logger := activity.GetLogger(ctx)

	fmt.Println(time.Now(), "checkDriverStatus activity -> Start")

	driverStatus := "ready"

	fmt.Println(time.Now(), "checkDriverStatus activity -> Finished")

	return driverStatus, nil
}
