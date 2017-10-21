package SIPServer

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func ReadAPIRequest(r *http.Request, out interface{}) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, out)
	return err
}

