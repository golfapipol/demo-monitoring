package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

const url = "mongodb://localhost:27017"

type Request struct {
	UUID      int
	Time      time.Duration
	StartTime time.Time
	Status    int
}

func main() {
	// create log file
	file, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	// connect
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer session.Close()

	log.SetOutput(file)
	gin.SetMode(gin.ReleaseMode)
	route := gin.Default()
	route.GET("/health", func(c *gin.Context) {
		var request Request
		var status int
		data := gin.H{}
		start := time.Now()
		random := rand.Intn(100)
		if random%3 == 0 {
			status = 200
			data["data"] = "lucky"
		} else {
			status = 500
			data["error"] = "unlucky"
			log.Println("line: 53 UUID ", random, "[error] @", start.Format(time.UnixDate), data)
		}
		c.JSON(status, data)
		end := time.Since(start)
		request.UUID = random
		request.Status = status
		request.StartTime = start
		request.Time = end
		session.DB("monitoring").C("requests").Insert(&request)
	})
	route.Run() // default 8080
}
