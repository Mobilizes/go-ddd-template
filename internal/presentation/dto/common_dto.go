package dto

const (
	MESSAGE_FAILED_GET_DATA_FROM_BODY = "failed get data from body"
	MESSAGE_FAILED_PROCESS_REQUEST    = "failed proses request"
	MESSAGE_FAILED_DENIED_ACCESS      = "denied access"

	MESSAGE_SUCCESS_GET_DATA    = "success get data"
	MESSAGE_SUCCESS_CREATE_DATA = "success create data"
	MESSAGE_SUCCESS_DELETE_DATA = "success delete data"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	Limit     int   `json:"limit"`
	Page      int   `json:"page"`
	TotalData int64 `json:"totalData"`
	TotalPage int   `json:"totalPage"`
}

func BuildResponseSuccess(message string, data any) Response {
	return Response{
		Message: message,
		Data:    data,
	}
}

func BuildPaginatedResponseSuccess(message string, data any, meta *Meta) Response {
	return Response{
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}

func BuildResponseFailed(message string, error any) Response {
	return Response{
		Message: message,
		Error:   error,
	}
}
