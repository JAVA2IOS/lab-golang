package config

import (
	"net/http"
	"encoding/json"
	"fmt"
	"log"
)

const (
	TimeForamt_yyyyMMddhhmmss = "2006-01-02 15:04:05"
	TimeForamt_yyyy_MM_dd_hh_mm_ss = "2006-01-02-15-04-05"
	TimeForamt_yyyyMMdd = "2006/01/02"
	TimeForamt_yyyy_MM_dd = "2006-01-02"
	TimeForamt_hhmmss = "15:04:05"
	TimeForamt_hh_mm_ss = "15_04_05"
)

type JsonHandler struct {
	Success bool `json:"success"`
	Data interface{} `json:"data"`
	Code int `json:"code"`
	Message string `json:"message"`
}

func HandlerAsyncHttpError(w http.ResponseWriter, err chan error) {
	select {
	case newErr := <- err:
		if newErr != nil {
			handler := JsonHandler{Success:false, Code: 201, Message : newErr.Error()}
			jsonString, _ := json.Marshal(handler)
			log.Printf("错误原因: %v\n", jsonString)
			fmt.Fprintf(w, string(jsonString))
		}
	default:
	}
}

func HandlerHttpJson(w http.ResponseWriter, data interface{}) {
	handler := JsonHandler{Success:true, Code: http.StatusOK, Data: data}
	jsonString, _ := json.Marshal(handler)
	log.Printf("成功: %v \n", string(jsonString))
	fmt.Fprintf(w, string(jsonString))
}

func HandlerHttpError(w http.ResponseWriter, err error) {
	if err != nil {
		handler := JsonHandler{Success:false, Code: 201, Message : err.Error()}
		jsonString, _ := json.Marshal(handler)
		log.Printf("错误原因: %v\n", jsonString)
		fmt.Fprintf(w, string(jsonString))
	}
}