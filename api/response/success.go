package response

import (
	"encoding/json"
	"net/http"
)

func WriteSuccess(w http.ResponseWriter, payload interface{}) error {
	w.WriteHeader(http.StatusOK)

	if nil != payload {
		encoder := json.NewEncoder(w)
		return encoder.Encode(payload)
	}

	return nil
}
