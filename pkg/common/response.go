package common

type Response struct {
	Code StatusCode  `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResponse(code StatusCode, data interface{}) Response {
	return Response{code, code.GetMsg(), data}
}

func SuccessResponse(data interface{}) Response {
	return NewResponse(Success, data)
}
