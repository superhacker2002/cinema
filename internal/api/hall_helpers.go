package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Helper function to retrieve the hall ID from the request URL parameters
func GetIntParam(r *http.Request, paramName string) (int, error) {
	vars := mux.Vars(r)
	paramValue := vars[paramName]

	paramInt, err := strconv.Atoi(paramValue)
	if err != nil {
		return 0, fmt.Errorf("Invalid %s", paramName)
	}

	return paramInt, nil
}

func GetHallID(r *http.Request) (int, error) {
	hallID, err := GetIntParam(r, "hallID")
	if err != nil {
		return 0, err
	}

	return hallID, nil
}

// Helper function to write a JSON response with the specified status code
func WriteResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
