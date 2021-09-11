package cassandra

import (
	"log"
	"os"

	"github.com/arshabbir/cassandraclient/domain/dto"
	"github.com/arshabbir/cassandraclient/utils"
	"github.com/gocql/gocql"
)

type client struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}
type Client interface {
	Create(dto.Student) *utils.ApiError
	Read(int) ([]dto.Student, *utils.ApiError)
	Update(int, dto.Student) *utils.ApiError
	Delete(int) *utils.ApiError
}

func NewDBClient() Client {
	//Get the Environment variable "CASSANDRACLUSTER"

	clusterIP := os.Getenv("CLUSTERIP")

	log.Println("ClusterIP environment  : ", clusterIP)

	cluster := gocql.NewCluster(clusterIP)
	cluster.Keyspace = "student"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()

	if err != nil {
		log.Println("Cassandra Session Creation Error..", err.Error())
		return nil
	}

	//defer session.Close()
	return &client{cluster: cluster, session: session}

}

func (c *client) Create(st dto.Student) *utils.ApiError {

	//Form the insert query & execute it

	//insertQuery := fmt.Sprintf("INSERT INTO studentdetails(id, name, class, marks) values(?, ?, ?, ?);")

	log.Println("Executing the insert query")
	if err := c.session.Query("INSERT INTO studentdetails(id, name, class, marks) values(?, ?, ?, ?);", st.Id, st.Name, st.Class, st.Marks).Consistency(gocql.Quorum).Exec(); err != nil {
		log.Println("Insert query error")
		return &utils.ApiError{Status: 0, Message: "Insert query error"}
	}

	return nil
}

func (c *client) Read(id int) ([]dto.Student, *utils.ApiError) {

	//Q

	var name, class string
	var idd, marks int

	iter := c.session.Query("SELECT id, name, class ,marks from studentdetails where id=?", id).Consistency(gocql.Quorum).Iter()
	var students = make([]dto.Student, iter.NumRows())

	log.Println("Number rows : ", iter.NumRows())

	for i := 0; iter.Scan(&idd, &name, &class, &marks); {
		//students = append(students, dto.Student{Name: name, Marks: marks, Id: idd, Class: class})
		students[i] = dto.Student{Name: name, Marks: marks, Id: idd, Class: class}
		i++

	}

	if err := iter.Close(); err != nil {
		log.Fatal("Error closing the iterator")
		return nil, nil
	}

	return students, nil
}

func (c *client) Update(id int, st dto.Student) *utils.ApiError {

	//updateQuery := fmt.Sprintf("Update studentdetails set  name=?, class=?, marks=? where id=?;")

	if err := c.session.Query("Update studentdetails set name=?, class=?, marks=? where id=? ;", st.Name, st.Class, st.Marks, id).Exec(); err != nil {
		log.Println("Update query error", err)
		return &utils.ApiError{Status: 0, Message: "Update query error"}
	}

	return nil
}

func (c *client) Delete(id int) *utils.ApiError {

	//deleteQuery := fmt.Sprintf("Delete studentdetails  where id=?;")

	if err := c.session.Query("Delete from studentdetails  where  id=? ;", id).Exec(); err != nil {
		log.Println("Delete query error", err)
		return &utils.ApiError{Status: 0, Message: "Delete query error"}
	}
	return nil
}
