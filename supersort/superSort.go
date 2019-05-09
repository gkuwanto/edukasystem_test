package supersort

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gkuwanto/edukasystem_test/logger"
	"github.com/julienschmidt/httprouter"
)

type userScore struct {
	ID    int     `json:"id_user"`
	Score float64 `json:"score"`
	Rank  int     `json:"rank"`
}

func SuperSort(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	go logger.LogAPICalls("SuperSort", r)
	db, err := sql.Open("mysql", "guest:codingtest@tcp(data.edukasystem.id:3306)/dummy")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	var ranks []userScore
	var user userScore

	rows, queryErr := db.Query("SELECT id_user, score FROM Score ORDER BY score DESC")
	if queryErr != nil {
		log.Fatal(queryErr)
	}
	defer rows.Close()
	numOfRows := 0
	for rows.Next() {
		scanErr := rows.Scan(&user.ID, &user.Score)
		if scanErr != nil {
			log.Fatal(scanErr)
		}
		ranks = append(ranks, user)
		ranks[numOfRows].Rank = numOfRows + 1
		numOfRows++
	}

	json.NewEncoder(w).Encode(ranks)
}
