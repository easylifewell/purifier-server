package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

type ReturnInfo struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func responseWithMessage(w http.ResponseWriter, message string, code int) {
	info := ReturnInfo{
		Message: message,
		Code:    code,
	}

	b, err := json.Marshal(info)
	if err != nil {
		log.Println(err)
	}
	ResponseWithJSON(w, b, code)
}

// Response200 return OK status
func Response200(w http.ResponseWriter, message string) {
	responseWithMessage(w, message, http.StatusOK)
}

// Response400 return when the request is bad
func Response400(w http.ResponseWriter, message string) {
	responseWithMessage(w, message, http.StatusBadRequest)
}

// Response500 return internal server error
func Response500(w http.ResponseWriter, message string) {
	responseWithMessage(w, message, http.StatusInternalServerError)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}
