package controller

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/arshabbir/consumer/utils"
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

	pprof.Register(sc.c)

	log.Println("Port environment  : ", port)
	sc.c.POST("/create", sc.Create)

	sc.c.GET("/metrics", prometheusHandler())
	sc.c.GET("/read/*id", sc.Read)

	sc.c.Run(port)

}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
