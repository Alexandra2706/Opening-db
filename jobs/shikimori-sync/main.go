package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	jsonURL       = "https://shikimori.one/api/animes/"
	logName       = "shikimori-log.log"
	maxNumberPage = 10
)

var log = logrus.New()

func main() {
	//file, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err == nil {
	//	log.Out = file
	//} else {
	//	log.Warn("Failed to log to file, using default stderr")
	//}

	for i := 1; i < maxNumberPage; i++ {
		MakeRequest(jsonURL + strconv.Itoa(i))
		time.Sleep(2 * time.Second)
	}

	//if _, err = file.WriteString("---------------"); err != nil {
	//	panic(err)
	//}
}

func init() {

	log.Out = os.Stdout

	log.SetFormatter(&logrus.JSONFormatter{})

	log.SetLevel(logrus.InfoLevel)
}

func MakeRequest(currentUrl string) {

	resp, err := http.Get(currentUrl)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != 200 {
		log.WithFields(logrus.Fields{
			"url":        currentUrl,
			"StatusCode": resp.StatusCode,
		}).Warn("Страница не найдена")
		//logrus.Warningf("%-10s: %-3d", currentUrl, resp.StatusCode)
	} else {
		log.WithFields(logrus.Fields{
			"url":        currentUrl,
			"StatusCode": resp.StatusCode,
		}).Info("Все OK")
	}

	//получить все содержимое страницы
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Println(string(body))
}
