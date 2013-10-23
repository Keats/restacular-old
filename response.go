package restacular

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	http.ResponseWriter
}

type ApiError struct {
	Code             int    `json:"code"`
	Status           int    `json:"status"`
	DeveloperMessage string `json:"developerMessage"`
	MoreInfo         string `json:"moreInfo"`
}

func (err *ApiError) Error() string {
	return err.DeveloperMessage
}

func NewInternalError() *ApiError {
	return &ApiError{500, 0, "Internal Server Error", ""}
}

func (resp *Response) Send(httpCode int, obj interface{}) {
	var content []byte

	content, err := json.Marshal(obj)
	if err != nil {
		resp.sendInternalError()
	}

	resp.writeResponse(httpCode, content)
}

func (resp *Response) SendError(apiError *ApiError) {
	var content []byte

	content, err := json.Marshal(apiError)
	if err != nil {
		resp.sendInternalError()
	}
	resp.writeResponse(apiError.Code, content)
}

func (resp *Response) sendInternalError() {
	resp.WriteHeader(http.StatusInternalServerError)
	resp.Write([]byte(""))
}

func (resp *Response) writeResponse(httpCode int, content []byte) {
	resp.Header().Set("content-type", "application/json")
	resp.WriteHeader(httpCode)
	resp.Write(content)
}
