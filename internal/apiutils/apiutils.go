package apiutils

import (
	"bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/cinemasession/service"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func IntPathParam(r *http.Request, varName string) (int, error) {
	vars := mux.Vars(r)
	varStr := vars[varName]
	varInt, err := strconv.Atoi(varStr)
	if err != nil {
		return 0, err
	}
	if varInt <= 0 {
		return 0, errors.New("parameter is less than zero")
	}
	return varInt, nil
}

func WriteResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println(err)
		http.Error(w, service.ErrInternalError.Error(), http.StatusInternalServerError)
	}
}

func WriteMsg(w http.ResponseWriter, msg string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Println(err)
		http.Error(w, service.ErrInternalError.Error(), http.StatusInternalServerError)
	}
}
