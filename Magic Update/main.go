package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type data struct {
	Msg  string `json:"message"`
	Data struct {
		ID   int `json:"id_user"`
		City int `json:"id_city"`
	}
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

func moveUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data1 := generateUserData()

	oldCityID := data1.Data.City
	newCityID := generateCityID()

	fmt.Println(oldCityID)
	fmt.Println(newCityID)

}

func main() {
	router := httprouter.New()
	router.GET("/moveuser", moveUser)

	log.Fatal(http.ListenAndServe(":8080", router))
}
