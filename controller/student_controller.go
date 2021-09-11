package controller

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/arshabbir/cassandraclient/utils"

	"github.com/arshabbir/cassandraclient/domain/dao"

	"github.com/arshabbir/cassandraclient/domain/dto"
	"github.com/gin-gonic/gin"
)

type studentController struct {
	c   *gin.Engine
	dao dao.StudentDAO
}

type StudentController interface {
	Create(c *gin.Context)
	Read(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Start()
}

func NewStudentController() StudentController {

	dao := dao.NewDAO()

	if dao == nil {
		log.Println("DAO creation error")
		return nil
	}

	log.Println("DAO creation successful .")

	return &studentController{c: gin.Default(), dao: dao}
}

func (sc *studentController) Create(c *gin.Context) {

	//Extract data from body

	var st dto.Student

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

func (sc *studentController) Read(c *gin.Context) {

	id, rerr := strconv.Atoi(c.Param("id"))

	if rerr != nil {
		log.Println("Error parsing the request")
		c.JSON(http.StatusOK, &utils.ApiError{Status: http.StatusBadRequest, Message: "Error parsing the request"})
		return
	}

	student, err := sc.dao.Read(id)
	if err != nil {

		log.Println("Error Reading data")
		c.JSON(http.StatusOK, &utils.ApiError{Status: http.StatusInternalServerError, Message: "Error Readin data"})
		return
	}

	//Marshal the result into json and send it

	c.JSON(http.StatusOK, &student)

	return
}

func (sc *studentController) Delete(c *gin.Context) {

	id, rerr := strconv.Atoi(c.Param("id"))

	if rerr != nil {
		log.Println("Error parsing the request")
		c.JSON(http.StatusOK, &utils.ApiError{Status: http.StatusBadRequest, Message: "Error parsing the request"})
		return
	}

	err := sc.dao.Delete(id)
	if err != nil {

		log.Println("Error Delete data")
		c.JSON(http.StatusOK, &utils.ApiError{Status: http.StatusInternalServerError, Message: "Error Deleting data"})
		return
	}

	c.JSON(http.StatusOK, &utils.ApiError{Status: http.StatusOK, Message: "Deletion Successful "})

	return
}

func (sc *studentController) Update(c *gin.Context) {

	//Extract data from body

	var st dto.Student

	if err := c.ShouldBindJSON(&st); err != nil {
		log.Println("Error parsing the request")
		c.JSON(http.StatusOK, &utils.ApiError{Status: http.StatusBadRequest, Message: "Error parsing the request"})
		return
	}

	if err := sc.dao.Update(st.Id, st); err != nil {

		log.Println("Error Updating into dao")
		c.JSON(http.StatusOK, &utils.ApiError{Status: http.StatusInternalServerError, Message: "Error Updating into dao"})
		return
	}

	log.Println("Insertion  successful .")
	c.JSON(http.StatusOK, &utils.ApiError{Status: http.StatusOK, Message: " Updating successful"})

	return
}

func (sc *studentController) Start() {

	port := os.Getenv("PORT")

	log.Println("Port environment  : ", port)
	sc.c.POST("/create", sc.Create)
	sc.c.GET("/delete/:id", sc.Delete)
	sc.c.PUT("/update", sc.Update)
	sc.c.GET("/read/:id", sc.Read)

	sc.c.Run(port)

}
