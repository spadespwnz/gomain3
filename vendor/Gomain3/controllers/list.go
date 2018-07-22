package controllers

import (
	"fmt"
	"github.com/globalsign/mgo"
	"os"
)

type (
	ListController struct {
		DB_NAME string
		Db      *mgo.Session
	}
)

func NewListController() *ListController {
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
	return &ListController{Db: sess, DB_NAME: DB_NAME}
}
