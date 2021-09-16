package common

import "context"

type Response struct {
	Code    StatusCode  `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	TraceId string      `json:"trace_id"`
}

func NewResponse(c context.Context, code StatusCode, data interface{}) Response {
	return Response{code, code.GetMsg(), data, c.Value("trace_id").(string)}
}

func SuccessResponse(c context.Context, data interface{}) Response {
	return NewResponse(c, Success, data)
}
