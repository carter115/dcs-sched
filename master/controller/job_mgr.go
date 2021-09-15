package controller

import (
	"context"
	"dcs-sched/master/config"
	common2 "dcs-sched/pkg/common"
	"encoding/json"
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

func (m *jobMgr) Save(c *gin.Context, job *common2.Job) error {
	jobKey := fmt.Sprintf("%s%d", common2.JOB_SAVE_DIR, job.Id)

	bs, err := json.Marshal(job)
	if err != nil {
		return err
	}
	resp, err := m.kv.Put(context.Background(), jobKey, string(bs))
	if err != nil {
		return err
	}
	gslog.Info(c, "etcd put response", resp)
	return nil
}
