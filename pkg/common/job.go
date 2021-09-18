package common

import (
	"context"
	"dcs-sched/pkg/db"
	"encoding/json"
	"github.com/carter115/gslog"
	"strconv"
	"strings"
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

// NewJob 生成任务
func NewJob(name, command, cronExpr string) *Job {
	return &Job{Id: NewJobId(), Name: name, Command: command, CronExpr: cronExpr}
}

// UnpackJob 反序列化Job
func UnpackJob(bs []byte) (*Job, error) {
	job := Job{}
	if err := json.Unmarshal(bs, &job); err != nil {
		gslog.Warning(context.Background(), "Unmarshal Job Error: Source:", string(bs))
		return nil, nil
	}
	return &job, nil
}

// ExtractJobId 从etcd的key中提取jobId
func ExtractJobId(jobKey string) int64 {
	s := strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
	id, err := strconv.Atoi(s)
	if err != nil {
		gslog.Warning(context.Background(), "ExtractJobId", err)
		return 0
	}
	return int64(id)
}

var (
	once      sync.Once
	defaultId int64 = 0
)

// NewJobId 生成jobId
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
