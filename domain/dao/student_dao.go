package dao

import (
	"log"

	"github.com/arshabbir/consumer/client/cassandra"
	"github.com/arshabbir/consumer/domain/dto"
	"github.com/arshabbir/consumer/utils"
)

type empDao struct {
	dbclient cassandra.Client
}

type EmpDAO interface {
	Create(dto.Emp) *utils.ApiError
	Read(int) ([]dto.Emp, *utils.ApiError)
}

func NewDAO() EmpDAO {
	dbclient := cassandra.NewDBClient()

	if dbclient == nil {
		log.Println("Error Creating the DAO ")
		return nil
	}

	return &empDao{dbclient: dbclient}
}

func (c *empDao) Create(st dto.Emp) *utils.ApiError {

	return c.dbclient.Create(st)
}

func (c *empDao) Read(id int) ([]dto.Emp, *utils.ApiError) {

	return c.dbclient.Read(id)
}
