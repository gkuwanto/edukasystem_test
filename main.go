package main

import (
	"log"
	"net/http"

	"github.com/gkuwanto/edukasystem_test/logger"
	"github.com/gkuwanto/edukasystem_test/magicupdate"
	"github.com/gkuwanto/edukasystem_test/supersort"
	"github.com/julienschmidt/httprouter"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	go logger.CleanLog()
	router := httprouter.New()
	router.GET("/MagicUpdate", magicupdate.MoveUser)
	router.GET("/SuperSorting", supersort.SuperSort)
	log.Fatal(http.ListenAndServe(":8000", router))
}
