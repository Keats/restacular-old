package restacular

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Request struct {
	http.Request
}

func (req *Request) DecodeJsonPayload(v interface{}) error {
	content, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, v)
	if err != nil {
		return err
	}
	return nil
}
