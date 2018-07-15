package main

import (
	"Gomain3/controllers"
	"log"
	"net/http"
)

func main() {
	apiController := controllers.NewApiController()
	defer apiController.Db.Close()
	http.HandleFunc("/", apiController.W)
	http.HandleFunc("/get", apiController.GetData)
	http.HandleFunc("/post", apiController.PostData)
	http.HandleFunc("/list", apiController.GetList)
	http.HandleFunc("/add/", apiController.AddList)
	http.HandleFunc("/drop/", apiController.DropDb)
	http.HandleFunc("/AddJpWord", apiController.AddJpWord)
	http.HandleFunc("/FindAll", apiController.FindAll)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
