package main

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

// PingHandler provides an api route for server health check
func PingHandler(w http.ResponseWriter, r *http.Request) {

	var data map[string]interface{}
	data = make(map[string]interface{})
	data["status"] = "success"
	result := make(map[string]interface{})
	result["result"] = "pong"
	result["registered"] = startTime.UTC()
	result["uptime"] = time.Since(startTime).Seconds()
	result["num_cores"] = runtime.NumCPU()
	data["data"] = result

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
