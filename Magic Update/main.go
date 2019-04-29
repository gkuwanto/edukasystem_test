package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func moveUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func main() {
	router := httprouter.New()
	router.GET("/moveuser", moveUser)

	log.Fatal(http.ListenAndServe(":8080", router))
}
