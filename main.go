package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type data struct {
	Msg  string `json:"message"`
	Data struct {
		ID   int `json:"id_user"`
		City int `json:"id_city"`
	}
}

type apiResponse struct {
	ID          int    `json:"id_user"`
	OldCityID   int    `json:"old_id_city"`
	OldCityName string `json:"old_city_name"`
	NewCityID   int    `json:"new_id_city"`
	NewCityName string `json:"new_city_name"`
}

func cleanLog() {
	for range time.Tick(time.Second * 15) {
		ioutil.WriteFile("log.txt", []byte(""), 0600)
	}
}

func logAPICalls(apiGateway string, request *http.Request) {
	userAgent := request.Header.Get("User-Agent")

	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	fmt.Fprintf(f, "%s %s %s \n", apiGateway, time.Now().Format("01/02/2006 15:04:05"), userAgent)
}

func generateUserData() data {
	url := "https://api.edukasystem.id/dummy/user/"
	data1 := data{}
	id := 300000000
	for data1.Msg != "ok" {
		rand.Seed(time.Now().UnixNano())
		id = 10000 + rand.Intn(10000)

		edukaClient := http.Client{
			Timeout: time.Second * 20,
		}

		req, err := http.NewRequest(http.MethodGet, url+strconv.Itoa(id), nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("User-Agent", "edukasystem-web-api-challange")

		res, getErr := edukaClient.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		jsonErr := json.Unmarshal(body, &data1)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		fmt.Println(data1.Msg)
	}
	return data1
}

func generateCityID() int {
	data1 := data{}
	id := 0
	for data1.Msg != "ok" {
		rand.Seed(time.Now().UnixNano())
		id = rand.Intn(1000)
		url := "https://api.edukasystem.id/dummy/city/"

		edukaClient := http.Client{
			Timeout: time.Second * 20,
		}
		req, err := http.NewRequest(http.MethodGet, url+strconv.Itoa(id), nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("User-Agent", "edukasystem-web-api-challange")

		res, getErr := edukaClient.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		jsonErr := json.Unmarshal(body, &data1)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
	}
	return id
}

type city struct {
	Msg  string `json:"message"`
	Data struct {
		ID   int    `json:"id_city"`
		Name string `json:"city_name"`
	}
}

func getCityName(id int) string {
	cityurl := "http://api.edukasystem.id/dummy/city/"
	idString := strconv.Itoa(id)
	cityClient := http.Client{
		Timeout: time.Second * 20,
	}
	req, err := http.NewRequest(http.MethodGet, cityurl+idString, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "edukasystem-web-api-challange")

	res, getErr := cityClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	city1 := city{}

	jsonErr := json.Unmarshal(body, &city1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return city1.Data.Name
}

func changeAPI(user int, city int) {
	url := "https://api.edukasystem.id/dummy/user/city"

	var jsonStr = []byte(
		`{
			"id_user":` + strconv.Itoa(user) + `,
			"id_city":` + strconv.Itoa(city) + `
		}`)

	changeClient := http.Client{
		Timeout: time.Second * 20,
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("User-Agent", "eduka")
	req.Header.Set("Content-Type", "application/json")
	resp, err := changeClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func moveUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	go logAPICalls("MagicUpdate", r)
	data1 := generateUserData()

	oldCityID := data1.Data.City
	newCityID := generateCityID()

	go changeAPI(data1.Data.ID, newCityID)
	apires := apiResponse{
		ID:          data1.Data.ID,
		OldCityID:   oldCityID,
		OldCityName: getCityName(oldCityID),
		NewCityID:   newCityID,
		NewCityName: getCityName(newCityID),
	}
	json.NewEncoder(w).Encode(apires)

}

type userScore struct {
	ID    int     `json:"id_user"`
	Score float64 `json:"score"`
	Rank  int     `json:"rank"`
}

type users [5000]userScore

func superSort(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	go logAPICalls("SuperSort", r)
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
		ranks[numOfRows].Rank = numOfRows
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

func main() {
	go cleanLog()
	router := httprouter.New()
	router.GET("/MagicUpdate", moveUser)
	router.GET("/SuperSorting", superSort)
	log.Fatal(http.ListenAndServe(":8080", router))
}
