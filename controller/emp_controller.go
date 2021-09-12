package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/arshabbir/consumer/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/arshabbir/consumer/domain/dao"

	"github.com/arshabbir/consumer/domain/dto"
)

type empController struct {
	c   *gin.Engine
	dao dao.EmpDAO
}

type EmpController interface {
	Create(c *gin.Context)
	Read(c *gin.Context)

	Consumer(chan int)

	Start()
}

func NewEmpController() EmpController {

	dao := dao.NewDAO()

	if dao == nil {
		log.Println("DAO creation error")
		return nil
	}

	log.Println("DAO creation successful .")

	return &empController{c: gin.Default(), dao: dao}
}

func (sc *empController) Create(c *gin.Context) {

	//Extract data from body

	var st dto.Emp

	if err := c.ShouldBindJSON(&st); err != nil {
		log.Println("Error parsing the request")
		c.JSON(http.StatusOK, &utils.ApiError{Status: http.StatusBadRequest, Message: "Error parsing the request"})
		return
	}

	if err := sc.dao.Create(st); err != nil {

		log.Println("Error inserting into dao")
		c.JSON(http.StatusOK, &utils.ApiError{Status: http.StatusInternalServerError, Message: "Error inserting into dao"})
		return
	}

	log.Println("Insertion  successful .")
	c.JSON(http.StatusOK, &utils.ApiError{Status: http.StatusOK, Message: "Insertion  successful "})

	return
}

func (sc *empController) Read(c *gin.Context) {

	id := c.Param("id")
	id = strings.TrimPrefix(id, "/")

	var emp []dto.Emp
	var err *utils.ApiError

	if id == "" {
		emp, err = sc.dao.ReadAll()

	} else {
		emp, err = sc.dao.Read(id)
	}
	if err != nil {

		log.Println("Error Reading data")
		c.JSON(http.StatusOK, &utils.ApiError{Status: http.StatusInternalServerError, Message: "Error Readin data"})
		return
	}

	//Marshal the result into json and send it

	c.JSON(http.StatusOK, &emp)

	return
}

func (sc *empController) Start() {

	port := os.Getenv("PORT")
	wait := make(chan int)

	pprof.Register(sc.c)

	log.Println("Port environment  : ", port)
	//	sc.c.POST("/create", sc.Create)

	//Start the Kakfa Consumer go routine

	go sc.Consumer(wait)

	sc.c.GET("/metrics", prometheusHandler())
	sc.c.GET("/v1/read/*id", sc.Read)

	sc.c.Run(port)

	<-wait

}

func (sc *empController) Consumer(wait chan int) {

	/*c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"auto.offset.reset": "earliest",
		"group.id":          "testgrouid",
	})*/

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_HOST"),
		"auto.offset.reset": "earliest",
		"group.id":          "testgrouid",
	})

	if err != nil {
		panic(err)
	}

	defer func() {
		wait <- 0
	}()
	//c.SubscribeTopics([]string{"myTopic1"}, nil)
	c.SubscribeTopics([]string{os.Getenv("TOPIC")}, nil)
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

			//Parse & Persist

			if err := sc.dao.Create(parseToken(string(msg.Value))); err != nil {

				log.Println("Error inserting into dao")
				continue

			}

			log.Println("Insertion  successful .")

		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)

		}

		//Persisit to DB
	}

	c.Close()

}

func parseToken(msg string) dto.Emp {

	//Name:name28162,Dept=OSS,EmplD:28162, Time=2021-09-11 15:42:38.125357173 +0000 UTC m=+26.014992376

	var emp dto.Emp

	tokens := strings.Split(msg, ",")

	//Parse Name :
	nTokens := strings.Split(tokens[0], ":")
	fmt.Printf("\nName Value : %s", nTokens[1])

	emp.Name = nTokens[1]

	//Parse Dept

	nTokens = strings.Split(tokens[1], "=")
	fmt.Printf("\nDept Value : %s", nTokens[1])

	emp.Dept = nTokens[1]

	nTokens = strings.Split(tokens[2], ":")
	fmt.Printf("\nEmp ID Value : %s", nTokens[1])
	emp.EmpID = nTokens[1]

	nTokens = strings.Split(tokens[3], "=")
	fmt.Printf("\n Timestamp Value : %s", nTokens[1])

	emp.TimeStamp = nTokens[1]

	return emp

}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
