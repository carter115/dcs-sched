package common

type StatusCode int

const (
	Success StatusCode = iota
	InvalidParam

	JobSaveError
	JobDeleteError
	JobListError
	JobKillError

	Unknown StatusCode = 9999
)

var statusCodeMap = map[StatusCode]string{
	Success:      "",
	InvalidParam: "invalid params",

	JobSaveError:   "保存任务失败",
	JobDeleteError: "删除任务失败",
	JobListError:   "列出任务失败",
	JobKillError:   "杀死任务失败",

	Unknown: "unknown",
}

func (c StatusCode) GetMsg() string {
	if s, ok := statusCodeMap[c]; ok {
		return s
	}
	return statusCodeMap[Unknown]
}
