package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/gorilla/mux"
	"github.com/d-smith/xtstatus/workitems"
	"net/http"
	"fmt"
	"log"
)

var (
	port = kingpin.Flag("port", "server listener port").Required().Int()
)

func main() {
	kingpin.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/workitems/{workitem}", workitems.WorkItemsHandler)

	http.Handle("/", r)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
