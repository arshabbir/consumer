package dao

import (
	"log"

	"github.com/arshabbir/cassandraclient/client/cassandra"
	"github.com/arshabbir/cassandraclient/domain/dto"
	"github.com/arshabbir/cassandraclient/utils"
)

type studentDao struct {
	dbclient cassandra.Client
}

type StudentDAO interface {
	Create(dto.Student) *utils.ApiError
	Read(int) ([]dto.Student, *utils.ApiError)
	Update(int, dto.Student) *utils.ApiError
	Delete(int) *utils.ApiError
}

func NewDAO() StudentDAO {
	dbclient := cassandra.NewDBClient()

	if dbclient == nil {
		log.Println("Error Creating the DAO ")
		return nil
	}

	return &studentDao{dbclient: dbclient}
}

func (c *studentDao) Create(st dto.Student) *utils.ApiError {

	return c.dbclient.Create(st)
}

func (c *studentDao) Read(id int) ([]dto.Student, *utils.ApiError) {

	return c.dbclient.Read(id)
}

func (c *studentDao) Update(id int, st dto.Student) *utils.ApiError {

	return c.dbclient.Update(id, st)
}

func (c *studentDao) Delete(id int) *utils.ApiError {

	return c.dbclient.Delete(id)

}
