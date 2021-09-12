package cassandra

import (
	"log"
	"os"
	"time"

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
	Read(string) ([]dto.Emp, *utils.ApiError)
	ReadAll() ([]dto.Emp, *utils.ApiError)
}

func NewDBClient() Client {
	//Get the Environment variable "CASSANDRACLUSTER"

	var session *gocql.Session
	var err error
	clusterIP := os.Getenv("CLUSTERIP")

	log.Println("ClusterIP environment  : ", clusterIP)

	cluster := gocql.NewCluster(clusterIP)
	//cluster.Keyspace = "student"
	cluster.Consistency = gocql.Quorum
	cluster.CQLVersion = "3"

	for {

		session, err = cluster.CreateSession()

		if err != nil {
			log.Println("Error creating session", err)

			time.Sleep(time.Second * 2)
		} else {
			break
		}

	}

	// create keyspaces
	err = session.Query("CREATE KEYSPACE IF NOT EXISTS consumersevents WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1};").Exec()
	if err != nil {
		log.Println(err)
		return nil
	}

	// create table
	//Name:XXXXX,Dept=OSS,EmplD:1234, Time=21-7-2021 21:00:10
	err = session.Query("CREATE TABLE IF NOT EXISTS consumersevents.events (name text, dept text, empid text, timestamp text,PRIMARY KEY (empid));").Exec()
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
	if err := c.session.Query("INSERT INTO consumersevents.events(name, dept, empid, timestamp) values(?, ?, ?, ?);", st.Name, st.Dept, st.EmpID, st.TimeStamp).Consistency(gocql.Quorum).Exec(); err != nil {
		log.Println("Insert query error", err)
		return &utils.ApiError{Status: 0, Message: "Insert query error"}
	}

	return nil
}

func (c *client) Read(id string) ([]dto.Emp, *utils.ApiError) {

	//Q

	var name, dept, empid, timestamp string

	iter := c.session.Query("SELECT name, dept, empid ,timestamp from consumersevents.events where empid=?", id).Consistency(gocql.Quorum).Iter()
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

func (c *client) ReadAll() ([]dto.Emp, *utils.ApiError) {

	//Q

	var name, dept, empid, timestamp string

	iter := c.session.Query("SELECT name, dept, empid ,timestamp from consumersevents.events ").Consistency(gocql.Quorum).Iter()
	var students = make([]dto.Emp, iter.NumRows())

	log.Println("Number rows : ", iter.NumRows())

	for i := 0; iter.Scan(&name, &dept, &empid, &timestamp); {
		//students = append(students, dto.Student{Name: name, Marks: marks, Id: idd, Class: class})
		students[i] = dto.Emp{Name: name, EmpID: empid, Dept: dept, TimeStamp: timestamp}
		i++

	}

	if err := iter.Close(); err != nil {
		log.Println("Error closing the iterator")
		return nil, nil
	}

	return students, nil
}
