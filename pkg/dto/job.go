package dto

type JobInput struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cron_expr"`
}

type JobOutput struct {
	Name string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cron_expr"`
}
