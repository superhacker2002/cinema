package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Helper function to retrieve the hall ID from the request URL parameters
func GetHallID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	hallID, err := strconv.Atoi(vars["hallID"])
	if err != nil {
		return 0, fmt.Errorf("Invalid hall ID")
	}
	return hallID, nil
}

// Helper function to write a JSON response with the specified status code
func WriteResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Helper function to handle errors and write an appropriate response
func HandleError(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, err.Error(), statusCode)
}
