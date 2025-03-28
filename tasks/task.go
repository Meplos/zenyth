package tasks

type TaskDef struct {
	Name string `json:"name"`
	Exec string `json:"exec"`
	Cron string `json:"cron"`
}
