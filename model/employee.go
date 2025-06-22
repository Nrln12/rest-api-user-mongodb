package model

type Employee struct {
	EmployeeId string `json:"employeeId,omitempty" bson:"employee_id"`
	Name       string `json:"name,omitempty" bson:"name"`
	Department string `json:"department,omitempty" bson:"department"`
}
