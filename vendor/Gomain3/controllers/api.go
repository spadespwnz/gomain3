package controllers

import (
	"Gomain3/models"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"io/ioutil"
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
	forceNew := false
	if r.Method != "GET" && r.Method != "POST" {
		SendError(w)
		return
	}
	r.ParseForm()
	romaji := r.Form.Get("romaji")
	if romaji == "" {
		fmt.Printf("Invalid word\n")
		fmt.Fprintf(w, "Invalid word")
		return
	}
	word := models.JPWord{
		State:  "new",
		Romaji: romaji,
	}
	meaning := r.Form.Get("meaning")
	if meaning != "" {
		word.Meaning = meaning
	}
	if r.Form.Get("force-new") == "true" {
		forceNew = true
	}

	word.Type = "JPWord"
	word.State = "new"

	col := api.Db.DB(api.DB_NAME).C("JP_COL")

	var existingWords []models.JPWord
	err := col.Find(bson.M{"romaji": romaji}).All(&existingWords)
	if err != nil {
		response := struct {
			Error    int    `json:"error"`
			ErrorMsg string `json:"error_message"`
		}{
			201,
			"Error Searching for Word in DB",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response)
		return
	}
	if len(existingWords) > 0 && !forceNew {
		response := struct {
			Error    int             `json:"error"`
			ErrorMsg string          `json:"error_message"`
			Words    []models.JPWord `json:"words"`
		}{
			100,
			"Word Already Exists in DB, need force-new",
			existingWords,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := col.Insert(word); err != nil {
		response := struct {
			Error    int    `json:"error"`
			ErrorMsg string `json:"error_message"`
		}{
			201,
			"Error Inserting Word in DB",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response)
		return
	}
	response := struct {
		Error    int           `json:"error"`
		ErrorMsg string        `json:"error_message"`
		Word     models.JPWord `json:"word"`
	}{
		0,
		"",
		word,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
	return
}

func (api ApiController) EditWords(w http.ResponseWriter, r *http.Request) {

	f, err := os.Open("./pages/words.html")
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		fmt.Fprintf(w, "Error loading page.")
		return
	}
	reader := bufio.NewReader(f)
	pageContent, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		fmt.Fprintf(w, "Error loading page.")
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(pageContent)
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
