package pendulum

type Query struct {
	TripID       string   `json:"trip_id"`
	ServiceID    string   `json:"service_id"`
	Origin       string   `json:"origin"`
	Destination  string   `json:"destination"`
	Start        int64    `json:"start"`
	End          int64    `json:"end"`
	Staff        []string `json:"staff"`
	Manager      []string `json:"manager"`
	CurrCityTask int      `json:"curr_city_task", default:0`
	Cities       []*City
}

type City struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Etd    int64  `json:"etd"`
	Status string `json:"status"`
	Tasks  []*Task
}

type Task struct {
	Name     []string `json:"name"`
	AlertMin int      `json:"alert_min"`
}
