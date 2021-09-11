package dto

type Emp struct {

	//Name:XXXXX,Dept=OSS,EmplD:1234, Time=21-7-2021 21:00:10

	Name      string `json:"name"`
	EmpID     string `json:"empid"`
	Dept      string `json:"dept"`
	TimeStamp string `json:"timestamp"`
}
