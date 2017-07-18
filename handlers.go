package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/schollz/jsonstore"
)

type Job struct {
	Message string `json:"message,omitempty"`
	Id      string `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
}

var Queues = map[string]chan Job{}
var KeyStore *jsonstore.JSONStore

func init() {
	Queues = make(map[string](chan Job))

	// initialize KeyStore
	var err error
	KeyStore, err = jsonstore.Open("jobs.json.gz")
	if err != nil {
		KeyStore = new(jsonstore.JSONStore)
	}

	// Check if keys already in KeyStore and add to Message Queue(s)
	keys := KeyStore.Keys()
	log.Println("adding pending jobs", len(keys))
	for i := range keys {
		key := keys[i]
		var job Job
		err := KeyStore.Get(key, &job)
		if nil != err {
			panic(err)
		}
		createQueueIfNotExists(job.Type)
		Queues[job.Type] <- job
	}
}

func createQueueIfNotExists(name string) {
	if _, ok := Queues[name]; !ok {
		Queues[name] = make(chan Job, 10000000)
	}
}

func NewJob(message string, queue_name string) (Job, error) {
	jobId, err := NewUUID()
	return Job{Message: message, Id: jobId, Type: queue_name}, err
}

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func WriteResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func GetQueuesHandler(w http.ResponseWriter, r *http.Request) {
	WriteResponseHeaders(w)

	data := make(map[string]int)
	for i := range Queues {
		data[i] = len(Queues[i])
	}

	js, err := json.Marshal(Response{Status: "ok", Data: data})
	if nil != err {
		http.Error(w, `{"status": "error", "message": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	WriteResponseHeaders(w)

	vars := mux.Vars(r)
	if "" == vars["q"] {
		http.Error(w, `{"status": "error", "message": "Must specify queue"}`, http.StatusBadRequest)
		return
	}

	if _, ok := Queues[vars["q"]]; !ok {
		http.Error(w, `{"status": "error", "message": "Queue not found"}`, http.StatusBadRequest)
		return
	}

	if 0 != len(Queues[vars["q"]]) {
		job := <-Queues[vars["q"]]
		js, err := json.Marshal(Response{Status: "ok", Data: job})
		if nil != err {
			http.Error(w, `{"status": "error", "message": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(js)

		// Remove job from KeyStore
		KeyStore.Delete(job.Id)
		jsonstore.Save(KeyStore, "jobs.json.gz")

		log.Println("Job retrieved", job)
	} else {
		js, err := json.Marshal(Response{Status: "ok", Data: nil})
		if nil != err {
			http.Error(w, `{"status": "error", "message": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}

}

func AddMessageHandler(w http.ResponseWriter, r *http.Request) {
	WriteResponseHeaders(w)

	// Check url params
	vars := mux.Vars(r)
	if "" == vars["q"] {
		http.Error(w, `{"status": "error", "message": "Must specify queue"}`, http.StatusBadRequest)
		return
	}

	// Create message queue if not exists
	createQueueIfNotExists(vars["q"])

	message, err := GetRequestBody(r)
	if nil != err {
		http.Error(w, `{"status": "error", "message": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Create new job
	job, err := NewJob(string(message), vars["q"])
	if nil != err {
		http.Error(w, `{"status": "error", "message": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Add job to KeyStore
	KeyStore.Set(job.Id, job)
	jsonstore.Save(KeyStore, "jobs.json.gz")

	// Add job to message queue
	Queues[vars["q"]] <- job

	log.Println("Job recieved", job)

	js, err := json.Marshal(Response{Status: "ok", Data: Job{Id: job.Id}})
	if nil != err {
		http.Error(w, `{"status": "error", "message": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func GetRId() string {
	return newJobId(8)
}

func GetRequestBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	return body, err
}
