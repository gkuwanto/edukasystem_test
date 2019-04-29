package main

import (
	"bytes"
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

type apiResponse struct {
	ID          int    `json:"id_user"`
	OldCityID   int    `json:"old_id_city"`
	OldCityName string `json:"old_city_name"`
	NewCityID   int    `json:"new_id_city"`
	NewCityName string `json:"new_city_name"`
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

func main() {
	router := httprouter.New()
	router.GET("/moveuser", moveUser)

	log.Fatal(http.ListenAndServe(":8080", router))
}
