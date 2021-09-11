package cassandra

import (
	"log"
	"os"

	"github.com/arshabbir/consumer/domain/dto"
	"github.com/arshabbir/consumer/utils"
	"github.com/gocql/gocql"
)

type client struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}
type Client interface {
	Create(dto.Emp) *utils.ApiError
	Read(int) ([]dto.Emp, *utils.ApiError)
}

func NewDBClient() Client {
	//Get the Environment variable "CASSANDRACLUSTER"

	clusterIP := os.Getenv("CLUSTERIP")

	log.Println("ClusterIP environment  : ", clusterIP)

	cluster := gocql.NewCluster(clusterIP)
	//cluster.Keyspace = "student"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()

	if err != nil {
		log.Println("Error creating session")
		return nil
	}

	// create keyspaces
	err = session.Query("CREATE KEYSPACE IF NOT EXISTS consumer WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1};").Exec()
	if err != nil {
		log.Println(err)
		return nil
	}

	// create table
	//Name:XXXXX,Dept=OSS,EmplD:1234, Time=21-7-2021 21:00:10
	err = session.Query("CREATE TABLE IF NOT EXISTS consumer.events (name text, dept text, empid text, time text,PRIMARY KEY (name, empid));").Exec()
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Table Creation Successful")

	//defer session.Close()
	return &client{cluster: cluster, session: session}

}

func (c *client) Create(st dto.Emp) *utils.ApiError {

	//Form the insert query & execute it

	//insertQuery := fmt.Sprintf("INSERT INTO studentdetails(id, name, class, marks) values(?, ?, ?, ?);")

	log.Println("Executing the insert query")
	if err := c.session.Query("INSERT INTO studentdetails(Name, Dept, EmpId, TimeStamp) values(?, ?, ?, ?);", st.Name, st.Dept, st.EmpID, st.TimeStamp).Consistency(gocql.Quorum).Exec(); err != nil {
		log.Println("Insert query error")
		return &utils.ApiError{Status: 0, Message: "Insert query error"}
	}

	return nil
}

func (c *client) Read(id int) ([]dto.Emp, *utils.ApiError) {

	//Q

	var name, dept, empid, timestamp string

	iter := c.session.Query("SELECT id, name, class ,marks from studentdetails where id=?", id).Consistency(gocql.Quorum).Iter()
	var students = make([]dto.Emp, iter.NumRows())

	log.Println("Number rows : ", iter.NumRows())

	for i := 0; iter.Scan(&name, &dept, &empid, &timestamp); {
		//students = append(students, dto.Student{Name: name, Marks: marks, Id: idd, Class: class})
		students[i] = dto.Emp{Name: name, EmpID: empid, Dept: dept, TimeStamp: timestamp}
		i++

	}

	if err := iter.Close(); err != nil {
		log.Fatal("Error closing the iterator")
		return nil, nil
	}

	return students, nil
}
