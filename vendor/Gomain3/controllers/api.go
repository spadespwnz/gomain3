package controllers

import (
	"Gomain3/models"
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo"
	"net/http"
	"os"
)

type (
	ApiController struct {
		DB_NAME string
		Db      *mgo.Session
	}
)

func NewApiController() *ApiController {
	DB_NAME := os.Getenv("MONGODB_DBNAME")
	if DB_NAME == "" {
		DB_NAME = "JP_DATA"
	}
	var dbUrl string = ""
	if os.Getenv("MONGODB_URI") != "" {
		dbUrl = os.Getenv("MONGODB_URI")
	} else {
		dbUrl = "mongodb://127.0.0.1:27017"
	}
	sess, err := mgo.Dial(dbUrl)
	if err != nil {
		fmt.Printf("%s", err)
	}
	return &ApiController{Db: sess, DB_NAME: DB_NAME}
}
func (api ApiController) W(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		SendError(w)
		return
	}
	fmt.Fprintf(w, "Welcome")
}

func (api ApiController) GetData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		SendError(w)
		return
	}
	fmt.Println("GET params were:", r.URL.Query())
	//param = r.URL.Query().Get("param");
	//if param != ""{}
	res := models.Response{
		Message: "Test Message 1",
		Error:   0,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", res)
}

func (api ApiController) PostData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendError(w)
		return
	}
	//r.FormValue("TestValue");

	res := models.Response{
		Message: "Test Message 2",
		Error:   0,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", res)
}

func (api ApiController) GetList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		SendError(w)
		return
	}
	myLists := []models.ListMeta{}
	myLists = append(myLists, models.ListMeta{"List 1", "Some Id 1"})
	myLists = append(myLists, models.ListMeta{"List 2", "Some Id 2"})
	body := models.Lists{Lists: myLists}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	json.NewEncoder(w).Encode(body)
}

func (api ApiController) AddList(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Path[len("/add/"):]
	fmt.Fprintf(w, "%s", t)
	return
	if r.Method != "GET" {
		SendError(w)
		return
	}
	list := r.URL.Query().Get("list")
	if list != "" {
		fmt.Fprintf(w, "%s", list)
	} else {
		fmt.Fprintf(w, "Invalid list name")
	}
}

func (api ApiController) AddJpWord(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		SendError(w)
		return
	}
	fmt.Println("GET params were:", r.URL.Query())
	romaji := r.URL.Query().Get("romaji")
	if romaji == "" {
		fmt.Printf("Invalid word\n")
		fmt.Fprintf(w, "Invalid word")
		return
	}
	word := models.JPWord{
		State:  "new",
		Romaji: romaji,
	}
	meaning := r.URL.Query().Get("meaning")
	if meaning != "" {
		word.Meaning = meaning
	}

	word.Type = "JPWord"
	word.State = "new"

	col := api.Db.DB(api.DB_NAME).C("JP_COL")

	if err := col.Insert(word); err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(word)
}

func (api ApiController) FindAll(w http.ResponseWriter, r *http.Request) {
	col := api.Db.DB(api.DB_NAME).C("JP_COL")
	var words []models.JPWord
	err := col.Find(nil).All(&words)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(words)
}

func (api ApiController) DropDb(w http.ResponseWriter, r *http.Request) {
	err := api.Db.DB(api.DB_NAME).C("JP_COL").DropCollection()
	if err != nil {
		fmt.Fprintf(w, "Failed to drop db")
		return
	}
	fmt.Fprintf(w, "DB dropped")
}

func SendError(w http.ResponseWriter) {
	res := models.Response{
		Message: "Page Not Found",
		Error:   404,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", res)
}
