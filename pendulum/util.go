package pendulum

import (
	"fmt"
	"time"
)

func epochToCron(hopEpoch int64) string {
	//set timezone
	loc, _ := time.LoadLocation("UTC")

	dateTime := time.Unix(hopEpoch, 0)
	fmt.Println("Next cron runs at: ", dateTime)
	localTime := dateTime.In(loc)

	min := localTime.Minute()
	hour := localTime.Hour()
	dayofmonth := localTime.Day()
	month := int(localTime.Month())
	dayofweek := int(localTime.Weekday())
	// year := dateTime.Year()

	cronExpression := fmt.Sprintf("%v %v %v %v %v", min, hour, dayofmonth, month, dayofweek)
	return cronExpression
}
