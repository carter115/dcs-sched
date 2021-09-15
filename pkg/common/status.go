package common

type StatusCode int

const (
	Success StatusCode = iota
	InvalidParam
	Unknown StatusCode = 9999
)

var statusCodeMap = map[StatusCode]string{
	Success:      "",
	InvalidParam: "invalid params",
	Unknown:      "unknown",
}

func (c StatusCode) GetMsg() string {
	if s, ok := statusCodeMap[c]; ok {
		return s
	}
	return statusCodeMap[Unknown]
}
