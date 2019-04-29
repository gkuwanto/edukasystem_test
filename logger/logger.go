package logger

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func CleanLog() {
	// Clears the log.txt file
	for range time.Tick(time.Second * 15) {
		ioutil.WriteFile("log.txt", []byte(""), 0600)
	}
}

func LogAPICalls(apiGateway string, request *http.Request) {
	//Appends a log message to log.txt
	userAgent := request.Header.Get("User-Agent")

	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	fmt.Fprintf(f, "%s %s %s \n", apiGateway, time.Now().Format("01/02/2006 15:04:05"), userAgent)
}
