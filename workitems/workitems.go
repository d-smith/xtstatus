package workitems

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
)

func WorkItemsHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	workitem := vars["workitem"]
	log.Println(workitem)
	rw.Write([]byte("placeholder\n"))
}