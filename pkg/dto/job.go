package dto

type JobInput struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cron_expr"`
}

type JobOutput struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cron_expr"`
}
