package main 

import (
    "io"
    "net/http"
    "github.com/gorilla/mux"
     log "github.com/sirupsen/logrus"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
    log.Info("/healthz")
    w.Header().Set("Content-Type", "application/json")
    io.WriteString(w, `{"status": "success"}`)
}

func init() {
    log.SetFormatter(&log.TextFormatter{})
    log.SetReportCaller(true)
}

func main() {
    
    router := mux.NewRouter()
    router.HandleFunc("/healthz", Healthz).Methods("GET")
    
    log.Info("Staring server...") 
    http.ListenAndServe(":8000", router)
}
