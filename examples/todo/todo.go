package main 

import (
    "io"
    "strconv"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    "github.com/jinzhu/gorm"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

var db, _ = gorm.Open("mysql", "root:root@/todolist?charset=utf8&parseTime=True&loc=Local")

// TODO: UPDATE gorm and  go get -u gorm.io/gorm 

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

func CreateItem(w http.ResponseWriter, r *http.Request) {
    description := r.FormValue("description")
    todo := &TodoItemModel{Description: description, Completed: false}
    db.Create(&todo)
    result := db.Last(&todo)
    log.WithFields(log.Fields{"description": description}).Info("Add new TodoItem. Saving to database.")
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result.Value)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, _ := strconv.Atoi(vars["id"])

    err := GetItemByID(id)
    if err == false {
        w.Header().Set("Content-Type", "application/json")
        io.WriteString(w, `{"updated": false, "error": "Item not found"}`)
    } else {
        completed, _ := strconv.ParseBool(r.FormValue("completed"))
        log.WithFields(log.Fields{"Id": id, "Completed": completed}).Info("Updating TodoItem")
        todo := &TodoItemModel{}
        db.First(&todo, id)
        todo.Completed = completed
        db.Save(&todo)
        w.Header().Set("Content-Type", "application/json")
        io.WriteString(w, `{"updated": true}`)
    }
}

func GetItemByID(Id int) bool {
    todo := &TodoItemModel{}
    result := db.First(&todo, Id)
    if result.Error != nil{
        log.Warn("TodoItem not found in database")
        return false
    }
    return true
}

func GetCompletedItems(w http.ResponseWriter, r *http.Request) {
    completedTodoItems := GetTodoItems(true)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(completedTodoItems)
}
func GetIncompletedItems(w http.ResponseWriter, r *http.Request) {
    completedTodoItems := GetTodoItems(false)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(completedTodoItems)
}

func GetTodoItems(completed bool) interface{} {
    var todos []TodoItemModel
    TodoItems := db.Where("completed = ?", completed).Find(&todos).Value
    return TodoItems
}

func GetAllTodoItems() interface{} {
    todoItems := []TodoItemModel{}
    db.Find(&todoItems)
    return todoItems
}

func GetAllItems(w http.ResponseWriter, r *http.Request) {
    allTodoItems := GetAllItems()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(allTodoItems)
} 

func DeleteItem(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, _ := strconv.Atoi(vars["id"])

    err := GetItemByID(id)
    if err == false {
        w.Header().Set("Content-Type", "application/json")
        io.WriteString(w, `{"updated": false, "error": "Item not found"}`
    } else {
        log.WithFields(log.Fields{"Id": id}).Info("Deleting TodoItem")
        todo := &TodoItemModel{}
        db.First(&todo, id)
        db.Delete(&todo)
        w.Header().Set("Content-Type", "application/json")
        io.WriteString(w, `{"deleted": true}`)
    }
}

func init() {
    log.SetFormatter(&log.TextFormatter{})
    log.SetReportCaller(true)
}

func main() {
    defer db.Close()

    db.Debug().DropTableIfExists(&TodoItemModel{})
    db.Debug().AutoMigrate(&TodoItemModel{})

    router := mux.NewRouter()
    router.HandleFunc("/healthz", Healthz).Methods("GET")
    router.HandleFunc("/todo", CreateItem).Methods("POST")
    router.HandlerFunc("/todo/all", GetAllItems).Methods("GET")
    router.HandleFunc("/todo/{id}", UpdateItem).Methods("POST")
    router.HandleFunc("/todo/{id}", DeleteItem).Methods("DELETE")
    router.HandleFunc("/todo-completed", GetCompletedItems).Methods("GET")
    router.HandleFunc("/todo-incomplete", GetIncompletedItems).Methods("GET")

    log.Info("Staring server on port:8000 ...") 
    http.ListenAndServe(":8000", router)
}
