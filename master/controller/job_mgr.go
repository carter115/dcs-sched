package controller

import (
	"dcs-sched/master/config"
	"dcs-sched/pkg/common"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/carter115/gslog"
	"github.com/gin-gonic/gin"
	"go.etcd.io/etcd/client/v3"
)

var JobMgr *jobMgr

type jobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

// InitJobMgr 初始化任务管理器
func InitJobMgr() error {
	conf := clientv3.Config{
		Endpoints:   config.Config.Etcd.Endpoints,
		DialTimeout: config.Config.Etcd.DialTimeout,
	}
	cli, err := clientv3.New(conf)
	if err != nil {
		return err
	}

	kv := clientv3.NewKV(cli)
	lease := clientv3.NewLease(cli)
	JobMgr = &jobMgr{
		client: cli,
		kv:     kv,
		lease:  lease,
	}
	return nil
}

// Save 保存任务
func (m *jobMgr) Save(c *gin.Context, job *common.Job) error {
	jobKey := fmt.Sprintf("%s%d", common.JOB_SAVE_DIR, job.Id)

	bs, err := json.Marshal(job)
	if err != nil {
		return err
	}
	resp, err := m.kv.Put(c, jobKey, string(bs))
	if err != nil {
		return err
	}
	gslog.Infof(c, "save job: %+v", resp)
	return nil
}

// Delete 删除任务
func (m *jobMgr) Delete(c *gin.Context, jobId int64) (*common.Job, error) {
	jobKey := fmt.Sprintf("%s%d", common.JOB_SAVE_DIR, jobId)
	delResp, err := m.kv.Delete(c, jobKey, clientv3.WithPrevKV())
	if err != nil {
		return nil, err
	}

	if len(delResp.PrevKvs) == 0 {
		return nil, errors.New(fmt.Sprintf("jobId not found: %d", jobId))
	}

	oldJobObj := common.Job{}
	bs := delResp.PrevKvs[0].Value
	if err := json.Unmarshal(bs, &oldJobObj); err != nil {
		gslog.Warning(c, "json unmarshal:", err)
		err = nil // 忽略解析旧值的出错
		return nil, err
	}
	gslog.Infof(c, "delete job: %+v", oldJobObj)
	return &oldJobObj, nil
}

// List 列出任务
func (m *jobMgr) List(c *gin.Context) ([]*common.Job, error) {
	jobKey := common.JOB_SAVE_DIR
	getResp, err := m.kv.Get(c, jobKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	if getResp.Count == 0 {
		return nil, errors.New("job list is empty")
	}

	jobs := make([]*common.Job, 0)
	for _, kvpair := range getResp.Kvs {
		obj := common.Job{}
		if err := json.Unmarshal(kvpair.Value, &obj); err != nil {
			gslog.Warning(c, "json unmarshal:", err)
		}
		jobs = append(jobs, &obj)
	}

	gslog.Infof(c, "job list: %+v", jobs)
	return jobs, nil
}

// Kill 杀死任务
func (m *jobMgr) Kill(c *gin.Context, jobId int64) error {
	jobKey := fmt.Sprintf("%s%d", common.JOB_KILL_DIR, jobId)
	leaseGrantResp, err := m.lease.Grant(c, 2)
	if err != nil {
		return err
	}

	leaseId := leaseGrantResp.ID
	_, err = m.kv.Put(c, jobKey, "", clientv3.WithLease(leaseId))

	return err
}
