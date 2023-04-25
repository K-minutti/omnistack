package main 

import (
    "io"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    "github.com/jinzhu/gorm"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

var db, _ = gorm.Open("mysql", "root:root@/todolist?charset=utf8&parseTime=True&loc=Local")

type TodoItemModel struct {
    Id int `gorm:"primary_key"`
    Description string
    Completed bool
}

func Healthz(w http.ResponseWriter, r *http.Request) {
    log.Info("/healthz")
    w.Header().Set("Content-Type", "application/json")
    io.WriteString(w, `{"status": "success"}`)
}


func init() {
    log.SetFormatter(&log.TextFormatter{})
    log.SetReportCaller(true)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
    description := r.FormValue("description")
    todo := &TodoItemModel{Description: description, Completed: false}
    db.Create(&todo)
    result := db.Last(&todo)
    log.WithFields(log.Fields{"description": description}).Info("Add new TodoItem. Saving to database.")
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result.Value)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {}
func GetItemByID(Id int) bool {return true}
func GetCompletedItems(w http.ResponseWriter, r *http.Request) {}
func GetIncompletedItems(w http.ResponseWriter, r *http.Request) {}
// func GetTodoItems(completed bool) interface{} {}
func DeleteItem(w http.ResponseWriter, r *http.Request) {}




func main() {
    defer db.Close()

    db.Debug().DropTableIfExists(&TodoItemModel{})
    db.Debug().AutoMigrate(&TodoItemModel{})

    router := mux.NewRouter()
    router.HandleFunc("/healthz", Healthz).Methods("GET")
    router.HandleFunc("/todo", CreateItem).Methods("POST")
    router.HandleFunc("/todo/{id}", UpdateItem).Methods("POST")
    router.HandleFunc("/todo/{id}", DeleteItem).Methods("DELETE")
    router.HandleFunc("/todo-completed", GetCompletedItems).Methods("GET")
    router.HandleFunc("/todo-incomplete", GetIncompletedItems).Methods("GET")

    log.Info("Staring server on port:8000 ...") 
    http.ListenAndServe(":8000", router)
}
