package common

const (
	JOB_EVENT_SAVE = iota
	JOB_EVENT_DELETE
	JOB_EVENT_KILL
)

// JobEvent 变化事件
type JobEvent struct {
	EventType int //SAVE or DELETE
	Job       *Job
}

// NewJobEvent 任务变化事件有2种：1）更新任务 2）删除任务
func NewJobEvent(eventType int, job *Job) *JobEvent {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}
