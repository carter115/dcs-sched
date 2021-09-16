package common

import (
	"dcs-sched/pkg/db"
	"sync"
)

const (
	JOB_ID_KEY   = "cron:jobid"
	JOB_SAVE_DIR = "/cron/job/"
	JOB_KILL_DIR = "/cron/killer/"
)

type Job struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cron_expr"`
}

func NewJob(name, command, cronExpr string) *Job {
	return &Job{Id: NewJobId(), Name: name, Command: command, CronExpr: cronExpr}
}

var (
	once      sync.Once
	defaultId int64 = 0
)

// 生成jobId
func NewJobId() int64 {
	once.Do(setDefaultJobId)
	return db.RedisCli.Incr(JOB_ID_KEY).Val()
}

func setDefaultJobId() {
	id, err := db.RedisCli.Get(JOB_ID_KEY).Int64()
	if err != nil || id == 0 {
		newid := db.RedisCli.Set(JOB_ID_KEY, defaultId, -1)
		if newid.Err() != nil {
			panic(newid.Err().Error())
		}
	}
}
