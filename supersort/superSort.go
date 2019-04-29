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

func merge(left, right []userScore) []userScore {
	size := len(left) + len(right)
	i := 0
	j := 0
	slice := make([]userScore, size, size)

	for k := 0; k < size; k++ {
		if i > len(left)-1 && j <= len(right)-1 {
			slice[k] = right[j]
			j++
		} else if j > len(right)-1 && i <= len(left)-1 {
			slice[k] = left[i]
			i++
		} else if left[i].Score > right[j].Score {
			slice[k] = left[i]
			i++
		} else {
			slice[k] = right[j]
			j++
		}
	}
	return slice
}

func mergeSort(slice []userScore) []userScore {

	if len(slice) < 2 {
		return slice
	}
	mid := (len(slice)) / 2
	return merge(mergeSort(slice[:mid]), mergeSort(slice[mid:]))
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

	rows, queryErr := db.Query("SELECT id_user, score FROM Score")
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
		numOfRows++
	}
	result := mergeSort(ranks)
	for i := 0; i < numOfRows; i++ {
		result[i].Rank = i + 1
	}
	json.NewEncoder(w).Encode(result)
}
