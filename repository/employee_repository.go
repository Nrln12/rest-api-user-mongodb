package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"user-api/model"
)

type EmployeeRepository struct {
	MongoCollection *mongo.Collection
}

func (er *EmployeeRepository) InsertEmployee(emp *model.Employee) (interface{}, error) {
	result, err := er.MongoCollection.InsertOne(context.Background(), emp)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (er *EmployeeRepository) FindEmployeeById(employeeId string) (*model.Employee, error) {
	var employee model.Employee
	err := er.MongoCollection.FindOne(context.Background(),
		bson.D{{Key: "employee_id", Value: employeeId}}).Decode(&employee)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (er *EmployeeRepository) FindAllEmployee() ([]model.Employee, error) {
	cur, err := er.MongoCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	var employees []model.Employee
	if err = cur.All(context.Background(), &employees); err != nil {
		return nil, fmt.Errorf("results decode error %s", err.Error())
	}
	return employees, nil
}

func (er *EmployeeRepository) UpdateEmployeeById(employeeId string, employee *model.Employee) (int64, error) {
	result, err := er.MongoCollection.UpdateOne(context.Background(),
		bson.D{{Key: "employee_id", Value: employeeId}},
		bson.D{{Key: "$set", Value: employee}})
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func (er *EmployeeRepository) DeleteEmployeeById(employeeId string) (int64, error) {
	result, err := er.MongoCollection.DeleteOne(context.Background(),
		bson.D{{Key: "employee_id", Value: employeeId}})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (er *EmployeeRepository) DeleteAllEmployee() (int64, error) {
	result, err := er.MongoCollection.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
