package v1

import "net/http"

type restSuccess struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func newCreateSuccess(data any) restSuccess {
	return restSuccess{
		Code:    http.StatusCreated,
		Data:    data,
		Message: "success create",
	}
}

func newGetSuccess(data any) restSuccess {
	return restSuccess{
		Code:    http.StatusOK,
		Data:    data,
		Message: "success get",
	}
}

func newUpdateSuccess(data any) restSuccess {
	return restSuccess{
		Code:    http.StatusOK,
		Data:    data,
		Message: "success update",
	}
}

func newDeleteSuccess() restSuccess {
	return restSuccess{
		Code:    http.StatusOK,
		Message: "success delete",
	}
}
