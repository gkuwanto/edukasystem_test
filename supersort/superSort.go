package supersort

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gkuwanto/edukasystem_web_api_golang/logger"
	"github.com/julienschmidt/httprouter"
)

type userScore struct {
	ID    int     `json:"id_user"`
	Score float64 `json:"score"`
	Rank  int     `json:"rank"`
}

type users [5000]userScore

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
	ranks := users{}

	rows, queryErr := db.Query("SELECT id_user, score FROM Score")
	if queryErr != nil {
		log.Fatal(queryErr)
	}
	defer rows.Close()
	numOfRows := 0
	for rows.Next() {
		scanErr := rows.Scan(&ranks[numOfRows].ID, &ranks[numOfRows].Score)
		if scanErr != nil {
			log.Fatal(scanErr)
		}
		ranks[numOfRows].Rank = numOfRows + 1
		numOfRows++
	}
	for i := 0; i < 5000; i++ {
		for j := 0; j < 5000; j++ {
			if ranks[i].Score > ranks[j].Score {
				temp := ranks[i].Rank
				ranks[i].Rank = ranks[j].Rank
				ranks[j].Rank = temp

				tem := ranks[i]
				ranks[i] = ranks[j]
				ranks[j] = tem
			}
		}
	}
	json.NewEncoder(w).Encode(ranks)
}
