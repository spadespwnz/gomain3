package main

import (
	"Gomain3/DiscoBot"
	"Gomain3/controllers"
	"log"
	"net/http"
	"os"
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
	http.HandleFunc("/edit", apiController.EditWords)
	http.HandleFunc("/update_word", apiController.ChangeJpWord)
	http.HandleFunc("/remove_word", apiController.RemoveJpWord)
	http.HandleFunc("/random_word", apiController.GetRandomWord)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
