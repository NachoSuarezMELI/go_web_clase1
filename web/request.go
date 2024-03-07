package web

import (
	"encoding/json"
	"errors"
	"net/http"
)

func RequestJsonProduct(r *http.Request, v interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("Content-Type header is not application/json")
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return nil
}
