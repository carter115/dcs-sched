package controller

import (
	"context"
	"dcs-sched/pkg/common"
	"dcs-sched/worker/config"
	"github.com/carter115/gslog"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var JobMgr *jobMgr

type jobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
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
	watcher := clientv3.NewWatcher(cli)
	JobMgr = &jobMgr{
		client:  cli,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}

	if err := JobMgr.watchJob(); err != nil {
		return err
	}
	if err := JobMgr.watchKiller(); err != nil {
		return err
	}

	return nil
}

// watchJob 监听任务
func (m *jobMgr) watchJob() error {
	ctx := context.Background()
	getResp, err := m.kv.Get(ctx, common.JOB_SAVE_DIR, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kvpair := range getResp.Kvs {
		if job := common.UnpackJob(kvpair.Value); job != nil {
			// TODO 把任务同步到scheudler
			gslog.Warning(ctx, "把任务同步到scheudler", job)
		}
	}

	// 从该revision向后监听变化事件
	go func() {
		startRevision := getResp.Header.Revision + 1
		watchChan := m.watcher.Watch(ctx, common.JOB_SAVE_DIR, clientv3.WithPrefix(), clientv3.WithRev(startRevision))

		// 处理监听事件
		for watchResp := range watchChan {
			for _, watchEvent := range watchResp.Events {
				gslog.Info(ctx, "Watch Event:", watchEvent)
				var jobEvent *common.JobEvent
				switch watchEvent.Type {
				case mvccpb.PUT:
					if job := common.UnpackJob(watchEvent.Kv.Value); job != nil {
						jobEvent = common.NewJobEvent(common.JOB_EVENT_SAVE, job)
					}

				case mvccpb.DELETE:
					job := &common.Job{Id: common.ExtractJobId(string(watchEvent.Kv.Key))}
					jobEvent = common.NewJobEvent(common.JOB_EVENT_DELETE, job)
				}
				// TODO 把(更新事件/删除事件)推给scheduler
				gslog.Warning(ctx, "把(更新事件/删除事件)推给scheduler", jobEvent)
			}
		}
	}()

	return nil
}

// watchKiller 监听杀死任务
func (m *jobMgr) watchKiller() error {
	ctx := context.Background()
	watchChan := m.watcher.Watch(ctx, common.JOB_KILL_DIR, clientv3.WithPrefix())
	for watchResp := range watchChan {
		for _, watchEvent := range watchResp.Events {
			switch watchEvent.Type {
			case mvccpb.PUT:
				job := &common.Job{Id: common.ExtractKillerId(string(watchEvent.Kv.Key))}
				jobEvent := common.NewJobEvent(common.JOB_EVENT_KILL, job)
				// TODO 把事件推给scheduler
				gslog.Warning(ctx, "把事件推给scheduler", jobEvent)

			case mvccpb.DELETE: // killer标记过期, 被自动删除
			}
		}
	}
	return nil
}
