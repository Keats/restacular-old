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

func (resp *Response) WriteError(apiError ApiError) {
	resp.WriteResponse(apiError.Code, apiError)
}

func (resp *Response) WriteResponse(httpCode int, obj interface{}) {
	resp.Header().Set("content-type", "application/json")

	var content []byte

	content, err := json.Marshal(obj)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(""))
	} else {
		resp.WriteHeader(httpCode)
		resp.Write(content)
	}
}
