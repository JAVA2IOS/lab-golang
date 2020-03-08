package config

import (
	"net/http"
	"encoding/json"
	"fmt"
)

type JsonHandler struct {
	Success bool `json:"success"`
	Data string `json:"data"`
	Code int `json:code`
	Err string `json:errorMessage`
}

func HandlerHttpError(w http.ResponseWriter, err chan error) {
	select {
	case newErr := <- err:
		if newErr != nil {
			handler := JsonHandler{Success:false, Code: 201, Err : newErr.Error()}
			jsonString, _ := json.Marshal(handler)
			fmt.Fprintf(w, string(jsonString))
		}
	default:
	}
}